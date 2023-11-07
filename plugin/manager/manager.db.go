package manager

type welcome struct {
	GrpID int64  `db:"gid"`
	Msg   string `db:"msg"`
}

type member struct {
	QQ int64 `db:"qq"`
	// GitHub username
	Ghun string `db:"ghun"`
}
