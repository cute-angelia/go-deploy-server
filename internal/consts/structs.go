package consts

// CmdMessage 发送给 client 执行命令
type CmdMessage struct {
	Action      string `json:"action"` // update rollback showlog
	Type        string `json:"type"`   // svn git
	Path        string `json:"path"`
	Reversion   string `json:"reversion,omitempty"`
	BeforDeploy string `json:"befor_deploy"`
	AfterDeploy string `json:"after_deploy"`
}
