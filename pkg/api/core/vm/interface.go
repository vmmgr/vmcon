package vm

import (
	"github.com/gorilla/websocket"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/controller/pkg/api/core"
	imaConStorage "github.com/vmmgr/imacon/pkg/api/core/storage"
	nodeCloudInit "github.com/vmmgr/node/pkg/api/core/tool/cloudinit"
	nodeVM "github.com/vmmgr/node/pkg/api/core/vm"
	"net/http"
	"time"
)

const (
	ID             = 0
	NodeID         = 1
	GroupID        = 2
	UUID           = 3
	NodeAndVNCPort = 4
	Lock           = 5
	UpdateName     = 100
	UpdateNode     = 101
	UpdateGroup    = 102
	UpdateBoot     = 103
	UpdateInfo     = 104
	UpdateAll      = 110
)

// channel定義(websocketで使用)
var Clients = make(map[*WebSocket]bool)
var Broadcast = make(chan WebSocketResult)
var ClientBroadcast = make(chan WebSocketResult)

// websocket用
type WebSocketResult struct {
	NodeID      uint      `json:"node_id"`
	Name        string    `json:"name"`
	Err         string    `json:"error"`
	CreatedAt   time.Time `json:"created_at"`
	UserToken   string    `json:"user_token"`
	AccessToken string    `json:"access_token"`
	UUID        string    `json:"uuid"`
	UserUUID    string    `json:"user_uuid"`
	VCPU        uint      `json:"vcpu"`
	Memory      uint      `json:"memory"`
	Status      int       `json:"status"`
	Code        uint      `json:"code"`
	GroupID     uint      `json:"group_id"`
	FilePath    string    `json:"file_path"`
	Admin       bool      `json:"admin"`
	Message     string    `json:"message"`
	Progress    uint      `json:"progress"`
}

type WebSocket struct {
	GroupID uint
	UserID  uint
	UUID    string
	Admin   bool
	Error   error
	Socket  *websocket.Conn
}

type Input struct {
	NodeID   uint   `json:"node_id"`
	GroupID  uint   `json:"group_id"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	OS       uint   `json:"os"` //32bit=> 32 64bit=> 64
	CPU      uint   `json:"cpu"`
	CPUModel string `json:"cpu_mode"`
	Memory   uint   `json:"memory"`
	VNCPort  uint   `json:"vnc_port"`
	Boot     uint   `json:"boot"` //0: hd 1:cdrom 2:floppy
	Lock     *bool  `json:"lock"`
}

type CreateAdmin struct {
	VM            nodeVM.VirtualMachine `json:"vm"`
	NodeID        uint                  `json:"node_id"`
	TemplateApply bool                  `json:"template_apply"`
	Template      Template              `json:"template"`
}

type DeleteAdmin struct {
	DeleteStorage bool `json:"delete_storage"`
}

type VirtualMachineStop struct {
	Force bool `json:"force"`
}

type GetImaCon struct {
	Status int               `json:"status"`
	Data   imaConStorage.Get `json:"data"`
}

type VMAll struct {
	VM      core.VM
	Storage core.Storage
	NIC     core.NIC
}

type Detail struct {
	VM   libvirtxml.Domain `json:"vm"`
	Stat uint              `json:"stat"`
	Node uint              `json:"node"`
}

type Template struct {
	Name            string   `json:"name"`
	Password        string   `json:"password"`
	NodeID          uint     `json:"node_id"`
	TemplateID      uint     `json:"template_id"`
	TemplatePlanID  uint     `json:"template_plan_id"`
	StorageCapacity uint     `json:"storage_capacity"`
	StoragePathType uint     `json:"storage_path_type"`
	IP              string   `json:"ip"`
	NetMask         string   `json:"netmask"`
	Gateway         string   `json:"gateway"`
	DNS             string   `json:"dns"`
	PCI             []string `json:"pci"`
	USB             []string `json:"usb"`
	NICType         string   `json:"nic_type"` //0:default 1~:custom
}

type CloudInit struct {
	HostName string                            `json:"hostname"`
	Name     string                            `json:"name"`
	Password string                            `json:"password"`
	Network  nodeCloudInit.NetworkConfigSubnet `json:"network"`
}

type Result struct {
	Status bool      `json:"status"`
	Error  string    `json:"error"`
	VM     []core.VM `json:"vm"`
}

type ResultAdmin struct {
	Status int      `json:"status"`
	Error  string   `json:"error"`
	VM     []Detail `json:"vm"`
}

type ResultOneAdmin struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	VM     Detail `json:"vm"`
}

type ResultOne struct {
	Status bool    `json:"status"`
	Error  string  `json:"error"`
	VM     core.VM `json:"vm"`
}

type ResultDatabase struct {
	Err error
	VMs []core.VM
}

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
