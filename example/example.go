package main

import (
	"fmt"
	"os"

	"github.com/yuin/goldmark"

	obsidian "github.com/powerman/goldmark-obsidian"
)

func main() {
	source := []byte(`
- [ ] Happy ==New== Year 📅 2026-01-01 ^first-task
- [x] Happy ~~Old~~ Year 📅 2025-01-01
	`)

	md := goldmark.New(
		goldmark.WithExtensions(
			obsidian.NewPlugTasks(),
			obsidian.NewObsidian(),
		),
	)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// <ul class="contains-task-list">
	// <li data-task="" class="task-list-item" id="^first-task"><input disabled="" type="checkbox" class="task-list-item-checkbox"> Happy <mark>New</mark> Year 📅 2026-01-01</li>
	// <li data-task="x" class="task-list-item is-checked"><input checked="" disabled="" type="checkbox" class="task-list-item-checkbox"> Happy <del>Old</del> Year 📅 2025-01-01</li>
	// </ul>
}
