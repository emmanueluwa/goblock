package network

type GetBlocksMessage struct {
	From uint32
	// if To is 0, max blocks will be returned
	To uint32
}

type GetStatusMessage struct {
}

type StatusMessage struct {
	//id of server
	ID            string
	Version       uint32
	CurrentHeight uint32
}
