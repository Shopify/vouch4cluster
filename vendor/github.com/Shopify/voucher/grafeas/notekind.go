package grafeas

import (
	"github.com/Shopify/voucher"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/common"
)

// getNoteKind translates the Voucher MetadataType into a Google Container Analysis NoteKind.
func getNoteKind(metadataType voucher.MetadataType) common.NoteKind {
	switch metadataType {
	case voucher.VulnerabilityType:
		return common.NoteKind_VULNERABILITY
	case voucher.BuildDetailsType:
		return common.NoteKind_BUILD
	case DiscoveryType:
		return common.NoteKind_DISCOVERY
	}
	return common.NoteKind_NOTE_KIND_UNSPECIFIED
}
