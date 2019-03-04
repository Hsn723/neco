package ipmi

// port from OpenIPMI

// Transport Network Function
const (
	IPMI_CMD_SET_LAN_CONFIG_PARMS =			0x01
	IPMI_CMD_GET_LAN_CONFIG_PARMS =			0x02
	IPMI_CMD_SUSPEND_BMC_ARPS =			0x03
	IPMI_CMD_GET_IP_UDP_RMCP_STATS =		0x04

	IPMI_CMD_SET_SERIAL_MODEM_CONFIG =		0x10
	IPMI_CMD_GET_SERIAL_MODEM_CONFIG =		0x11
	IPMI_CMD_SET_SERIAL_MODEM_MUX =			0x12
	IPMI_CMD_GET_TAP_RESPONSE_CODES =		0x13
	IPMI_CMD_SET_PPP_UDP_PROXY_XMIT_DATA =		0x14
	IPMI_CMD_GET_PPP_UDP_PROXY_XMIT_DATA =		0x15
	IPMI_CMD_SEND_PPP_UDP_PROXY_PACKET =		0x16
	IPMI_CMD_GET_PPP_UDP_PROXY_RECV_DATA =		0x17
	IPMI_CMD_SERIAL_MODEM_CONN_ACTIVE =		0x18
	IPMI_CMD_CALLBACK =				0x19
	IPMI_CMD_SET_USER_CALLBACK_OPTIONS =		0x1a
	IPMI_CMD_GET_USER_CALLBACK_OPTIONS =		0x1b

	IPMI_CMD_SOL_ACTIVATING =			0x20
	IPMI_CMD_SET_SOL_CONFIGURATION_PARAMETERS =	0x21
	IPMI_CMD_GET_SOL_CONFIGURATION_PARAMETERS =	0x22
)
