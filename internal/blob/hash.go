package blob

import "notgit/utils"

func hash(b *Blob) {
	b.Hash = utils.Hash("blob", b.Content)
}
