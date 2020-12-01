package peer


type PeerPicker interface {
	PickPeer (key string) string
}