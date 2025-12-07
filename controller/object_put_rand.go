package controller

import (
	"io"
	"math/rand"
	"path"
	"strconv"
	"time"

	"github.com/sj14/sss/util"
)

func (c *Controller) ObjectPutRand(dest string, size, count uint64, cfg ObjectPutConfig) error {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range count {
		limit := &io.LimitedReader{R: random, N: int64(size)}

		nameLen := max(3, len(strconv.Itoa(int(count))))
		objectName := util.RandomString(uint64(nameLen), util.LettersLower)
		fp := path.Join(dest, objectName)

		err := c.objectPut(limit, size, fp, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}
