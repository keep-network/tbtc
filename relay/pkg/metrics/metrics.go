package metrics

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/metrics"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
	"github.com/keep-network/tbtc/relay/pkg/node"
)

var logger = log.Logger("relay-metrics")

const (
	// DefaultChainMetricsTick is the default duration of the
	// observation tick for chain-related metrics.
	DefaultChainMetricsTick = 10 * time.Minute
	// DefaultNodeMetricsTick is the default duration of the
	// observation tick for node-related metrics.
	DefaultNodeMetricsTick = 10 * time.Second
)

// Initialize set up the metrics registry and enables metrics server.
func Initialize(
	port int,
) (*metrics.Registry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := metrics.NewRegistry()

	registry.EnableServer(port)

	return registry, true
}

// ObserveBtcChainConnectivity triggers an observation process of the
// btc_chain_connectivity metric.
func ObserveBtcChainConnectivity(
	ctx context.Context,
	registry *metrics.Registry,
	btcHandle btc.Handle,
	tick time.Duration,
) {
	input := func() float64 {
		_, err := btcHandle.GetBlockCount()

		if err != nil {
			return 0
		}

		return 1
	}

	observe(
		ctx,
		"btc_chain_connectivity",
		input,
		registry,
		validateTick(tick, DefaultChainMetricsTick),
	)
}

// ObserveHostChainConnectivity triggers an observation process of the
// host_chain_connectivity metric.
func ObserveHostChainConnectivity(
	ctx context.Context,
	registry *metrics.Registry,
	hostChain chain.Handle,
	tick time.Duration,
) {
	input := func() float64 {
		_, err := hostChain.GetBestKnownDigest()

		if err != nil {
			return 0
		}

		return 1
	}

	observe(
		ctx,
		"host_chain_connectivity",
		input,
		registry,
		validateTick(tick, DefaultChainMetricsTick),
	)
}

// ObserveBlockForwarding triggers an observation process of the
// block_forwarding metric.
func ObserveBlockForwarding(
	ctx context.Context,
	registry *metrics.Registry,
	nodeStats node.Stats,
	tick time.Duration,
) {
	input := func() float64 {
		isEnabled := nodeStats.BlockForwardingEnabled()

		if !isEnabled {
			return 0
		}

		return 1
	}

	observe(
		ctx,
		"block_forwarding",
		input,
		registry,
		validateTick(tick, DefaultNodeMetricsTick),
	)
}

// ObserveBlockForwardingErrors triggers an observation process of the
// block_forwarding_errors metric.
func ObserveBlockForwardingErrors(
	ctx context.Context,
	registry *metrics.Registry,
	nodeStats node.Stats,
	tick time.Duration,
) {
	input := func() float64 {
		return float64(nodeStats.BlockForwardingErrors())
	}

	observe(
		ctx,
		"block_forwarding_errors",
		input,
		registry,
		validateTick(tick, DefaultNodeMetricsTick),
	)
}

// ObserveUniqueBlocksPulled triggers an observation process of the
// unique_blocks_pulled metric.
func ObserveUniqueBlocksPulled(
	ctx context.Context,
	registry *metrics.Registry,
	nodeStats node.Stats,
	tick time.Duration,
) {
	input := func() float64 {
		return float64(nodeStats.UniqueBlocksPulled())
	}

	observe(
		ctx,
		"unique_blocks_pulled",
		input,
		registry,
		validateTick(tick, DefaultNodeMetricsTick),
	)
}

// ObserveUniqueBlocksPushed triggers an observation process of the
// unique_blocks_pushed metric.
func ObserveUniqueBlocksPushed(
	ctx context.Context,
	registry *metrics.Registry,
	nodeStats node.Stats,
	tick time.Duration,
) {
	input := func() float64 {
		return float64(nodeStats.UniqueBlocksPushed())
	}

	observe(
		ctx,
		"unique_blocks_pushed",
		input,
		registry,
		validateTick(tick, DefaultNodeMetricsTick),
	)
}

func observe(
	ctx context.Context,
	name string,
	input metrics.ObserverInput,
	registry *metrics.Registry,
	tick time.Duration,
) {
	observer, err := registry.NewGaugeObserver(name, input)
	if err != nil {
		logger.Warnf("could not create gauge observer [%v]", name)
		return
	}

	observer.Observe(ctx, tick)
}

func validateTick(tick time.Duration, defaultTick time.Duration) time.Duration {
	if tick > 0 {
		return tick
	}

	return defaultTick
}
