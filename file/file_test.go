package file

import "testing"

func TestTreeDir(t *testing.T) {
	t.Logf("%+v", TreeDir("/Users/ezbuy/Projects/ezbuy/goflow/src/github.com/JodeZer/lazydog", -1))
	t.Logf("%+v", TreeDir("/Users/ezbuy/Projects/ezbuy/goflow/src/github.com/JodeZer/lazydog", 2))
}

func TestListFile(t *testing.T) {
	t.Logf("%+v", ListGoFileByPaths(TreeDir("/Users/ezbuy/Projects/ezbuy/goflow/src/github.com/JodeZer/lazydog", -1), false))
}

func TestJumperBackup(t *testing.T) {
	jp := &Jumper{}
	jp.BackupPath("../example")
}

func TestJumperRestore(t *testing.T) {
	jp := &Jumper{}
	jp.RestorePath("../example")
}
