package updater

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/neco"
	"github.com/cybozu-go/neco/storage"
	"github.com/cybozu-go/well"
)

// Server represents neco-updater server
type Server struct {
	session *concurrency.Session
	storage storage.Storage
	timeout time.Duration

	checker ReleaseChecker
}

// NewServer returns a Server
func NewServer(session *concurrency.Session, storage storage.Storage, timeout time.Duration) Server {
	return Server{
		session: session,
		storage: storage,
		timeout: timeout,

		checker: NewReleaseChecker(storage),
	}
}

// Run runs neco-updater
func (s Server) Run(ctx context.Context) error {

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	e := concurrency.NewElection(s.session, storage.KeyLeader)

RETRY:
	select {
	case <-s.session.Done():
		return errors.New("session has been orphaned")
	default:
	}

	err = e.Campaign(ctx, hostname)
	if err != nil {
		return err
	}
	leaderKey := e.Key()

	log.Info("I am the leader", map[string]interface{}{
		"session": s.session.Lease(),
	})

	env := well.NewEnvironment(ctx)
	env.Go(func(ctx context.Context) error {
		return s.runLoop(ctx, leaderKey)
	})
	env.Go(func(ctx context.Context) error {
		return s.checker.Run(ctx)
	})
	env.Stop()
	err = env.Wait()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	err2 := e.Resign(ctxWithTimeout)
	if err2 != nil {
		return err2
	}
	if err == storage.ErrNoLeader {
		log.Warn("lost the leadership", nil)
		goto RETRY
	}
	return err
}

func (s Server) runLoop(ctx context.Context, leaderKey string) error {
	var target string
	var current *neco.UpdateRequest
	var currentRev int64

	// Updater continues last update without create update reuqest with "skipRequest = true"
	var skipRequest bool

	current, currentRev, err := s.storage.GetRequestWithRev(ctx)
	if err == nil {
		target = current.Version
		if current.Stop {
			log.Info("Last updating is failed, wait for retrying", map[string]interface{}{
				"version": current.Version,
			})
			err := s.storage.WaitRequestDeletion(ctx, current, currentRev)
			if err != nil {
				return err
			}
		} else {
			log.Info("Last updating is still in progress, wait for workers", map[string]interface{}{
				"version": current.Version,
			})
			skipRequest = true
		}
	} else if err != nil && err != storage.ErrNotFound {
		return err
	}

	for {
		if len(target) == 0 {
			if s.checker.HasUpdate() {
				target = s.checker.GetLatest()
			}
		}
		if len(target) != 0 {
			// Found new update
			for {
				if !skipRequest {
					current, currentRev, err = s.startUpdate(ctx, target, leaderKey)
					if err == ErrNoMembers {
						break
					} else if err != nil {
						return err
					}
				}
				skipRequest = false

				aborted, err := s.storage.WaitWorkers(ctx, current, currentRev)
				if err == nil {
					s.notifySlackSucceeded(ctx, *current)
					break
				}
				if err == storage.ErrTimedOut {
					s.notifySlackTimeout(ctx, *current)
				}
				if !aborted {
					return err
				}
				s.notifySlackServerFailure(ctx, *current, err.Error())
				err = s.stopUpdate(ctx, current, leaderKey)
				if err != nil {
					return err
				}
				err = s.storage.WaitRequestDeletion(ctx, current, currentRev)
				if err != nil {
					return err
				}
			}
		}

		timeout, err := s.storage.WaitForMemberUpdated(ctx, current, currentRev)
		if timeout {
			// Clear target to check latest update
			target = ""
		} else if err != nil {
			return err
		}
	}
}

// startUpdate starts update with tag.  It returns created request and its
// modified revision.  It returns ErrNoMembers if no bootservers are registered
// in etcd.
func (s Server) startUpdate(ctx context.Context, tag, leaderKey string) (*neco.UpdateRequest, int64, error) {
	servers, err := s.storage.GetBootservers(ctx)
	if err != nil {
		return nil, 0, err
	}
	if len(servers) == 0 {
		log.Info("No bootservers exists in etcd", map[string]interface{}{})
		return nil, 0, ErrNoMembers
	}
	log.Info("Starting updating", map[string]interface{}{
		"version": tag,
		"servers": servers,
	})
	r := neco.UpdateRequest{
		Version:   tag,
		Servers:   servers,
		Stop:      false,
		StartedAt: time.Now(),
	}
	err = s.storage.PutRequest(ctx, r, leaderKey)
	if err != nil {
		return nil, 0, err
	}
	return s.storage.GetRequestWithRev(ctx)
}

// stopUpdate of the current request
func (s Server) stopUpdate(ctx context.Context, req *neco.UpdateRequest, leaderKey string) error {
	req.Stop = true
	return s.storage.PutRequest(ctx, *req, leaderKey)
}

func (s Server) notifySlackSucceeded(ctx context.Context, req neco.UpdateRequest) error {
	slack, err := s.newSlackClient(ctx)
	if err == storage.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}

	att := Attachment{
		Color:      ColorGood,
		AuthorName: "Boot server updater",
		Title:      "Update completed successfully",
		Text:       "Updating on boot servers are completed successfully :tada: :tada: :tada:",
		Fields: []AttachmentField{
			{Title: "Version", Value: req.Version, Short: true},
			{Title: "Servers", Value: fmt.Sprintf("%v", req.Servers), Short: true},
			{Title: "Started at", Value: req.StartedAt.Format(time.RFC3339), Short: true},
		},
	}
	payload := Payload{Attachments: []Attachment{att}}
	return slack.PostWebHook(ctx, payload)
}

func (s Server) notifySlackServerFailure(ctx context.Context, req neco.UpdateRequest, message string) error {
	slack, err := s.newSlackClient(ctx)
	if err == storage.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}

	att := Attachment{
		Color:      ColorDanger,
		AuthorName: "Boot server updater",
		Title:      "Failed to update boot servers",
		Text:       "Failed to update boot servers due to some worker return(s) error :crying_cat_face:.  Please fix it manually.",
		Fields: []AttachmentField{
			{Title: "Version", Value: req.Version, Short: true},
			{Title: "Servers", Value: fmt.Sprintf("%v", req.Servers), Short: true},
			{Title: "Started at", Value: req.StartedAt.Format(time.RFC3339), Short: true},
			{Title: "Reason", Value: message, Short: true},
		},
	}
	payload := Payload{Attachments: []Attachment{att}}
	return slack.PostWebHook(ctx, payload)
}

func (s Server) notifySlackTimeout(ctx context.Context, req neco.UpdateRequest) error {
	slack, err := s.newSlackClient(ctx)
	if err == storage.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}

	att := Attachment{
		Color:      ColorDanger,
		AuthorName: "Boot server updater",
		Title:      "Update failed on the boot servers",
		Text:       "Failed to update boot servers due to timed-out from worker updates :crying_cat_face:.  Please fix it manually.",
		Fields: []AttachmentField{
			{Title: "Version", Value: req.Version, Short: true},
			{Title: "Servers", Value: fmt.Sprintf("%v", req.Servers), Short: true},
			{Title: "Started at", Value: req.StartedAt.Format(time.RFC3339), Short: true},
		},
	}
	payload := Payload{Attachments: []Attachment{att}}
	return slack.PostWebHook(ctx, payload)
}

func (s Server) newSlackClient(ctx context.Context) (*SlackClient, error) {
	webhookURL, err := s.storage.GetSlackNotification(ctx)
	if err == storage.ErrNotFound {
		return nil, storage.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var http *http.Client

	proxyURL, err := s.storage.GetProxyConfig(ctx)
	if err == storage.ErrNotFound {
	} else if err != nil {
		return nil, err
	} else {
		u, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		http = newHTTPClient(u)
	}

	return &SlackClient{URL: webhookURL, HTTP: http}, nil
}
