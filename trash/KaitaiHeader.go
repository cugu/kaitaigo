type KaitaiHeader struct {
    _io     *Stream
    _parent interface{}
    _root   interface{}
}

func (k *KaitaiHeader) Init(io *Stream, parent interface{}, root interface{}) {
    k._io = io
    k._parent = parent
    k._root = root
}

func (k *KaitaiHeader) IO() *Stream {
    return k._io
}

func (k *KaitaiHeader) Parent() interface{} {
    return k._parent
}

func (k *KaitaiHeader) Root() interface{} {
    return k._root
}

func (k *KaitaiHeader) SetIO(value *Stream) {
    k._io = value
}

func (k *KaitaiHeader) SetParent(value interface{}) {
    k._parent = value
}

func (k *KaitaiHeader) SetRoot(value interface{}) {
    k._root = value
}
