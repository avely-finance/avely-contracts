package contracts

type Protocol struct {
	Zproxy *Zproxy
	Zimpl *Zimpl
	Aimpl *AZil
	Buffer *BufferContract
	Holder *HolderContract
}

func NewProtocol(zproxy *Zproxy, zimpl *Zimpl, azil *AZil, buffer *BufferContract, holder *HolderContract) *Protocol {
	return &Protocol{
		Zproxy: zproxy,
		Zimpl: zimpl,
		Aimpl: azil,
		Buffer: buffer,
		Holder: holder,
	}
}