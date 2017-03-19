package rapport

const (
	prefixKeyspaceByte = '\x01'
	prefixSegmentByte  = '\x02'
	prefixSystemByte   = '\x03'
)

var (
	// PrefixUserKey is the bytes that prefix user keyspace keys
	PrefixUserKey = []byte{prefixKeyspaceByte}

	// PrefixSegmentKey is the bytes that prefix user segment keys. Segment keys
	// are support data for the user keyspace
	PrefixSegmentKey = []byte{prefixSegmentByte}

	// PrefixSystemKey is the bytes that prefix system keyspace keys
	PrefixSystemKey = []byte{prefixSystemByte}
)
