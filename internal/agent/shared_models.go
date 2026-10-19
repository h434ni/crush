package agent

import (
	"context"

	"github.com/charmbracelet/crush/internal/config"
)

// ModelSet is the large and small model an agent sends requests through.
type ModelSet struct {
	Large, Small Model
}

// coordinatorModels is every model the coordinator's agents share. Sub-agents
// get their own large model because copilot bills them as agent-initiated;
// the small model is always built that way, so one serves everyone.
type coordinatorModels struct {
	large, subLarge, small Model
}

// rebuildModels rebuilds every shared model from the live config and
// publishes them together. Rebuilds are serialized, so one that read an older
// config can never publish over one that read a newer credential.
func (c *coordinator) rebuildModels(ctx context.Context) error {
	c.rebuildMu.Lock()
	defer c.rebuildMu.Unlock()

	large, err := c.buildModel(ctx, config.SelectedModelTypeLarge, false)
	if err != nil {
		return err
	}
	subLarge, err := c.buildModel(ctx, config.SelectedModelTypeLarge, true)
	if err != nil {
		return err
	}
	small, err := c.buildModel(ctx, config.SelectedModelTypeSmall, true)
	if err != nil {
		return err
	}
	c.models.Set(coordinatorModels{large: large, subLarge: subLarge, small: small})
	return nil
}

func (c *coordinator) mainModels() ModelSet {
	m := c.models.Get()
	return ModelSet{Large: m.large, Small: m.small}
}

func (c *coordinator) subAgentModels() ModelSet {
	m := c.models.Get()
	return ModelSet{Large: m.subLarge, Small: m.small}
}

// smallOnlyModels serves agents that use the small model for everything.
func (c *coordinator) smallOnlyModels() ModelSet {
	m := c.models.Get()
	return ModelSet{Large: m.small, Small: m.small}
}
