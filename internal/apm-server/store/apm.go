package store

import (
	"context"
)

type ApmStore interface {
	Get(ctx context.Context)
}

func Get() {

}
