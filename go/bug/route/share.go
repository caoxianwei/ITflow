package route

import (
	"github.com/hyahm/xmux"
	"itflow/bug/handle"
)

var Share *xmux.GroupRoute

func init() {
	Share = xmux.NewGroupRoute()
	Share.Pattern("/share/list").Get(handle.ShareList)
	Share.Pattern("/share/upload").Post(handle.ShareUpload)
	Share.Pattern("/share/mkdir").Post(handle.ShareMkdir)
	Share.Pattern("/share/remove").Get(handle.ShareRemove)
	Share.Pattern("/share/rename").Post(handle.ShareRename)
	//router.HandleFunc("/share/down", handle.ShareDownload)
	Share.Pattern("/share/down").Get(handle.ShareShow)
}
