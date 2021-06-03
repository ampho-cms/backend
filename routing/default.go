// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package routing

// NewDefault creates a new default router.
func NewDefault() Router {
	return NewGorillaMux()
}
