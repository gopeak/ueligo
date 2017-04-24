package web

type Root struct {
	Code int
	Msg  string
	Data interface{}
}

type ListType struct {
	Mine   map[string]string
	Friend FriendType
	Group  []map[string]string
}

type FriendType struct {
	Groupname string
	Online    int
	Id        int
	List      []map[string]string
}
