package structs

// Link struct of Fido's link document
type Link struct {
	ID              string `bson:"_id" json:"id"`
	OriginalPath    string `bson:"original_path" json:"original_path"`
	SharedPath      string `bson:"shared_path" json:"shared_path"`
	NodePath        string `bson:"node_path" json:"node_path"`
	ReplicationNode string `bson:"replication_node" json:"replication_node"`
	MIMEType        string `bson:"mime_type" json:"mime_type"`
	FileSize        int64  `bson:"file_size" json:"file_size"`
	CreatedAt       int64  `bson:"created_at" json:"created_at"`
}
