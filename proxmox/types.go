package proxmox

type VirtualMachine struct {
	Name       string `json:"name"`
	Vmid       int    `json:"vmid"`
	Status     string `json:"status"`
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
}

type ConfigData struct {
	Config Config `json:"data"`
}

type Config struct {
	Digest   string `json:"digest"`
	Name     string `json:"name"`
	Scsi0    string `json:"scsi0"`
	Cores    int    `json:"cores"`
	Ostype   string `json:"ostype"`
	Smbios1  string `json:"smbios1"`
	Scsihw   string `json:"scsihw"`
	Memory   int    `json:"memory"`
	Boot     string `json:"boot"`
	Net0     string `json:"net0"`
	Sockets  int    `json:"sockets"`
	Numa     int    `json:"numa"`
	Bootdisk string `json:"bootdisk"`
	Ide2     string `json:"ide2"`
}

type Qemu struct {
	Proxmox *Proxmox
	Vmid    int    `json:"vmid"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type QemuData struct {
	Qemus []Qemu `json:"data"`
}

type Ticket struct {
	Data TicketData `json:"data"`
}

type TicketData struct {
	Cap struct {
		Access struct {
			UserModify    int `json:"User.Modify"`
			GroupAllocate int `json:"Group.Allocate"`
		} `json:"access"`
		Dc struct {
			SysAudit int `json:"Sys.Audit"`
		} `json:"dc"`
		Vms struct {
			VMConfigOptions   int `json:"VM.Config.Options"`
			PermissionsModify int `json:"Permissions.Modify"`
			VMConfigHWType    int `json:"VM.Config.HWType"`
			VMPowerMgmt       int `json:"VM.PowerMgmt"`
			VMConfigCPU       int `json:"VM.Config.CPU"`
			VMBackup          int `json:"VM.Backup"`
			VMConfigDisk      int `json:"VM.Config.Disk"`
			VMConfigMemory    int `json:"VM.Config.Memory"`
			VMConsole         int `json:"VM.Console"`
			VMAudit           int `json:"VM.Audit"`
			VMMonitor         int `json:"VM.Monitor"`
			VMConfigNetwork   int `json:"VM.Config.Network"`
			VMAllocate        int `json:"VM.Allocate"`
			VMConfigCDROM     int `json:"VM.Config.CDROM"`
			VMSnapshot        int `json:"VM.Snapshot"`
			VMMigrate         int `json:"VM.Migrate"`
			VMClone           int `json:"VM.Clone"`
		} `json:"vms"`
		Nodes struct {
			SysModify    int `json:"Sys.Modify"`
			SysPowerMgmt int `json:"Sys.PowerMgmt"`
			SysConsole   int `json:"Sys.Console"`
			SysAudit     int `json:"Sys.Audit"`
			SysSyslog    int `json:"Sys.Syslog"`
		} `json:"nodes"`
		Storage struct {
			DatastoreAudit            int `json:"Datastore.Audit"`
			PermissionsModify         int `json:"Permissions.Modify"`
			DatastoreAllocate         int `json:"Datastore.Allocate"`
			DatastoreAllocateTemplate int `json:"Datastore.AllocateTemplate"`
			DatastoreAllocateSpace    int `json:"Datastore.AllocateSpace"`
		} `json:"storage"`
	} `json:"cap"`
	Username            string `json:"username"`
	CSRFPreventionToken string `json:"CSRFPreventionToken"`
	Ticket              string `json:"ticket"`
}
