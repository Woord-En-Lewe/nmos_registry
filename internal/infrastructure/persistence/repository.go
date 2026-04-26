package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/persistence/db"
	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
)

type SQLiteRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewSQLiteRepository(sqlDB *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		queries: db.New(sqlDB),
		db:      sqlDB,
	}
}

// Node operations
func (r *SQLiteRepository) UpsertNode(ctx context.Context, node registry.Node) error {
	return r.queries.UpsertNode(ctx, db.UpsertNodeParams{
		ID:          node.ID,
		ApiVersion:  node.ApiVersion,
		Version:     node.Version,
		Label:       node.Label,
		Description: node.Description,
		Tags:        nilToEmptyJSON(node.Tags),
		Href:        node.Href,
		Hostname:    toNullString(node.Hostname),
		Api:         nilToEmptyJSON(node.Api),
		Caps:        nilToEmptyJSON(node.Caps),
		Services:    nilToEmptyJSON(node.Services),
		Clocks:      nilToEmptyJSON(node.Clocks),
		Interfaces:  nilToEmptyJSON(node.Interfaces),
		LastSeen:    sql.NullTime{Time: time.Now(), Valid: true},
	})
}

func (r *SQLiteRepository) DeleteNode(ctx context.Context, id string) error {
	return r.queries.DeleteNode(ctx, id)
}

func (r *SQLiteRepository) GetNode(ctx context.Context, id string) (registry.Node, error) {
	n, err := r.queries.GetNode(ctx, id)
	if err != nil {
		return registry.Node{}, err
	}
	return registry.Node{
		ID:          n.ID,
		ApiVersion:  n.ApiVersion,
		Version:     n.Version,
		Label:       n.Label,
		Description: n.Description,
		Tags:        n.Tags,
		Href:        n.Href,
		Hostname:    fromNullString(n.Hostname),
		Api:         n.Api,
		Caps:        n.Caps,
		Services:    n.Services,
		Clocks:      n.Clocks,
		Interfaces:  n.Interfaces,
		LastSeen:    n.LastSeen.Time,
	}, nil
}

func (r *SQLiteRepository) ListNodes(ctx context.Context) ([]registry.Node, error) {
	dbNodes, err := r.queries.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	nodes := make([]registry.Node, len(dbNodes))
	for i, n := range dbNodes {
		nodes[i] = registry.Node{
			ID:          n.ID,
			ApiVersion:  n.ApiVersion,
			Version:     n.Version,
			Label:       n.Label,
			Description: n.Description,
			Tags:        n.Tags,
			Href:        n.Href,
			Hostname:    fromNullString(n.Hostname),
			Api:         n.Api,
			Caps:        n.Caps,
			Services:    n.Services,
			Clocks:      n.Clocks,
			Interfaces:  n.Interfaces,
			LastSeen:    n.LastSeen.Time,
		}
	}
	return nodes, nil
}

func (r *SQLiteRepository) UpdateNodeHealth(ctx context.Context, id string) error {
	return r.queries.UpdateNodeLastSeen(ctx, db.UpdateNodeLastSeenParams{
		ID:       id,
		LastSeen: sql.NullTime{Time: time.Now(), Valid: true},
	})
}

func (r *SQLiteRepository) GetExpiredNodes(ctx context.Context, expirationTime time.Time) ([]string, error) {
	return r.queries.GetExpiredNodes(ctx, sql.NullTime{Time: expirationTime, Valid: true})
}

// Device operations
func (r *SQLiteRepository) UpsertDevice(ctx context.Context, device registry.Device) error {
	return r.queries.UpsertDevice(ctx, db.UpsertDeviceParams{
		ID:          device.ID,
		ApiVersion:  device.ApiVersion,
		NodeID:      device.NodeID,
		Version:     device.Version,
		Label:       device.Label,
		Description: device.Description,
		Tags:        device.Tags,
		Type:        device.Type,
		Senders:     device.Senders,
		Receivers:   device.Receivers,
		Controls:    device.Controls,
	})
}

func (r *SQLiteRepository) DeleteDevice(ctx context.Context, id string) error {
	return r.queries.DeleteDevice(ctx, id)
}

func (r *SQLiteRepository) GetDevice(ctx context.Context, id string) (registry.Device, error) {
	d, err := r.queries.GetDevice(ctx, id)
	if err != nil {
		return registry.Device{}, err
	}
	return registry.Device{
		ID:          d.ID,
		ApiVersion:  d.ApiVersion,
		NodeID:      d.NodeID,
		Version:     d.Version,
		Label:       d.Label,
		Description: d.Description,
		Tags:        d.Tags,
		Type:        d.Type,
		Senders:     d.Senders,
		Receivers:   d.Receivers,
		Controls:    d.Controls,
	}, nil
}

func (r *SQLiteRepository) ListDevices(ctx context.Context) ([]registry.Device, error) {
	dbDevices, err := r.queries.ListDevices(ctx)
	if err != nil {
		return nil, err
	}
	devices := make([]registry.Device, len(dbDevices))
	for i, d := range dbDevices {
		devices[i] = registry.Device{
			ID:          d.ID,
			ApiVersion:  d.ApiVersion,
			NodeID:      d.NodeID,
			Version:     d.Version,
			Label:       d.Label,
			Description: d.Description,
			Tags:        d.Tags,
			Type:        d.Type,
			Senders:     d.Senders,
			Receivers:   d.Receivers,
			Controls:    d.Controls,
		}
	}
	return devices, nil
}

// Source operations
func (r *SQLiteRepository) UpsertSource(ctx context.Context, source registry.Source) error {
	return r.queries.UpsertSource(ctx, db.UpsertSourceParams{
		ID:          source.ID,
		ApiVersion:  source.ApiVersion,
		DeviceID:    source.DeviceID,
		Version:     source.Version,
		Label:       source.Label,
		Description: source.Description,
		Tags:        source.Tags,
		GrainRate:   source.GrainRate,
		Format:      source.Format,
		Caps:        source.Caps,
		Parents:     source.Parents,
		ClockName:   toNullString(source.ClockName),
	})
}

func (r *SQLiteRepository) DeleteSource(ctx context.Context, id string) error {
	return r.queries.DeleteSource(ctx, id)
}

func (r *SQLiteRepository) GetSource(ctx context.Context, id string) (registry.Source, error) {
	s, err := r.queries.GetSource(ctx, id)
	if err != nil {
		return registry.Source{}, err
	}
	return registry.Source{
		ID:          s.ID,
		ApiVersion:  s.ApiVersion,
		DeviceID:    s.DeviceID,
		Version:     s.Version,
		Label:       s.Label,
		Description: s.Description,
		Tags:        s.Tags,
		GrainRate:   s.GrainRate,
		Format:      s.Format,
		Caps:        s.Caps,
		Parents:     s.Parents,
		ClockName:   fromNullString(s.ClockName),
	}, nil
}

func (r *SQLiteRepository) ListSources(ctx context.Context) ([]registry.Source, error) {
	dbSources, err := r.queries.ListSources(ctx)
	if err != nil {
		return nil, err
	}
	sources := make([]registry.Source, len(dbSources))
	for i, s := range dbSources {
		sources[i] = registry.Source{
			ID:          s.ID,
			ApiVersion:  s.ApiVersion,
			DeviceID:    s.DeviceID,
			Version:     s.Version,
			Label:       s.Label,
			Description: s.Description,
			Tags:        s.Tags,
			GrainRate:   s.GrainRate,
			Format:      s.Format,
			Caps:        s.Caps,
			Parents:     s.Parents,
			ClockName:   fromNullString(s.ClockName),
		}
	}
	return sources, nil
}

// Flow operations
func (r *SQLiteRepository) UpsertFlow(ctx context.Context, flow registry.Flow) error {
	return r.queries.UpsertFlow(ctx, db.UpsertFlowParams{
		ID:                     flow.ID,
		ApiVersion:             flow.ApiVersion,
		SourceID:               flow.SourceID,
		DeviceID:               flow.DeviceID,
		Version:                flow.Version,
		Label:                  flow.Label,
		Description:            flow.Description,
		Tags:                   flow.Tags,
		Format:                 flow.Format,
		MediaType:              toNullString(flow.MediaType),
		SampleRate:             flow.SampleRate,
		BitDepth:               toNullInt64(flow.BitDepth),
		DidSdid:                flow.DidSdid,
		GrainRate:              flow.GrainRate,
		FrameWidth:             toNullInt64(flow.FrameWidth),
		FrameHeight:            toNullInt64(flow.FrameHeight),
		InterlaceMode:          toNullString(flow.InterlaceMode),
		Colorspace:             toNullString(flow.Colorspace),
		Components:             flow.Components,
		TransferCharacteristic: toNullString(flow.TransferCharacteristic),
	})
}

func (r *SQLiteRepository) DeleteFlow(ctx context.Context, id string) error {
	return r.queries.DeleteFlow(ctx, id)
}

func (r *SQLiteRepository) GetFlow(ctx context.Context, id string) (registry.Flow, error) {
	f, err := r.queries.GetFlow(ctx, id)
	if err != nil {
		return registry.Flow{}, err
	}
	return registry.Flow{
		ID:                     f.ID,
		ApiVersion:             f.ApiVersion,
		SourceID:               f.SourceID,
		DeviceID:               f.DeviceID,
		Version:                f.Version,
		Label:                  f.Label,
		Description:            f.Description,
		Tags:                   f.Tags,
		Format:                 f.Format,
		MediaType:              fromNullString(f.MediaType),
		SampleRate:             f.SampleRate,
		BitDepth:               fromNullInt64ToInt(f.BitDepth),
		DidSdid:                f.DidSdid,
		GrainRate:              f.GrainRate,
		FrameWidth:             fromNullInt64ToInt(f.FrameWidth),
		FrameHeight:            fromNullInt64ToInt(f.FrameHeight),
		InterlaceMode:          fromNullString(f.InterlaceMode),
		Colorspace:             fromNullString(f.Colorspace),
		Components:             f.Components,
		TransferCharacteristic: fromNullString(f.TransferCharacteristic),
	}, nil
}

func (r *SQLiteRepository) ListFlows(ctx context.Context) ([]registry.Flow, error) {
	dbFlows, err := r.queries.ListFlows(ctx)
	if err != nil {
		return nil, err
	}
	flows := make([]registry.Flow, len(dbFlows))
	for i, f := range dbFlows {
		flows[i] = registry.Flow{
			ID:                     f.ID,
			ApiVersion:             f.ApiVersion,
			SourceID:               f.SourceID,
			DeviceID:               f.DeviceID,
			Version:                f.Version,
			Label:                  f.Label,
			Description:            f.Description,
			Tags:                   f.Tags,
			Format:                 f.Format,
			MediaType:              fromNullString(f.MediaType),
			SampleRate:             f.SampleRate,
			BitDepth:               fromNullInt64ToInt(f.BitDepth),
			DidSdid:                f.DidSdid,
			GrainRate:              f.GrainRate,
			FrameWidth:             fromNullInt64ToInt(f.FrameWidth),
			FrameHeight:            fromNullInt64ToInt(f.FrameHeight),
			InterlaceMode:          fromNullString(f.InterlaceMode),
			Colorspace:             fromNullString(f.Colorspace),
			Components:             f.Components,
			TransferCharacteristic: fromNullString(f.TransferCharacteristic),
		}
	}
	return flows, nil
}

// Senders
func (r *SQLiteRepository) UpsertSender(ctx context.Context, sender registry.Sender) error {
	return r.queries.UpsertSender(ctx, db.UpsertSenderParams{
		ID:                     sender.ID,
		ApiVersion:             sender.ApiVersion,
		DeviceID:               sender.DeviceID,
		FlowID:                 toNullString(sender.FlowID),
		Version:                sender.Version,
		Label:                  sender.Label,
		Description:            sender.Description,
		Tags:                   sender.Tags,
		Transport:              sender.Transport,
		ManifestHref:           toNullString(sender.ManifestHref),
		InterfaceBindings:      sender.InterfaceBindings,
		SubscriptionReceiverID: toNullString(sender.SubscriptionReceiverID),
		SubscriptionActive:     toNullBool(sender.SubscriptionActive),
	})
}

func (r *SQLiteRepository) DeleteSender(ctx context.Context, id string) error {
	return r.queries.DeleteSender(ctx, id)
}

func (r *SQLiteRepository) GetSender(ctx context.Context, id string) (registry.Sender, error) {
	s, err := r.queries.GetSender(ctx, id)
	if err != nil {
		return registry.Sender{}, err
	}
	return registry.Sender{
		ID:                     s.ID,
		ApiVersion:             s.ApiVersion,
		DeviceID:               s.DeviceID,
		FlowID:                 fromNullString(s.FlowID),
		Version:                s.Version,
		Label:                  s.Label,
		Description:            s.Description,
		Tags:                   s.Tags,
		Transport:              s.Transport,
		ManifestHref:           fromNullString(s.ManifestHref),
		InterfaceBindings:      s.InterfaceBindings,
		SubscriptionReceiverID: fromNullString(s.SubscriptionReceiverID),
		SubscriptionActive:     fromNullBool(s.SubscriptionActive),
	}, nil
}

func (r *SQLiteRepository) ListSenders(ctx context.Context) ([]registry.Sender, error) {
	dbSenders, err := r.queries.ListSenders(ctx)
	if err != nil {
		return nil, err
	}
	senders := make([]registry.Sender, len(dbSenders))
	for i, s := range dbSenders {
		senders[i] = registry.Sender{
			ID:                     s.ID,
			ApiVersion:             s.ApiVersion,
			DeviceID:               s.DeviceID,
			FlowID:                 fromNullString(s.FlowID),
			Version:                s.Version,
			Label:                  s.Label,
			Description:            s.Description,
			Tags:                   s.Tags,
			Transport:              s.Transport,
			ManifestHref:           fromNullString(s.ManifestHref),
			InterfaceBindings:      s.InterfaceBindings,
			SubscriptionReceiverID: fromNullString(s.SubscriptionReceiverID),
			SubscriptionActive:     fromNullBool(s.SubscriptionActive),
		}
	}
	return senders, nil
}

// Receivers
func (r *SQLiteRepository) UpsertReceiver(ctx context.Context, receiver registry.Receiver) error {
	return r.queries.UpsertReceiver(ctx, db.UpsertReceiverParams{
		ID:                   receiver.ID,
		ApiVersion:           receiver.ApiVersion,
		DeviceID:             receiver.DeviceID,
		Version:              receiver.Version,
		Label:                receiver.Label,
		Description:          receiver.Description,
		Tags:                 receiver.Tags,
		Transport:            receiver.Transport,
		Format:               receiver.Format,
		Caps:                 receiver.Caps,
		InterfaceBindings:    receiver.InterfaceBindings,
		SubscriptionSenderID: toNullString(receiver.SubscriptionSenderID),
		SubscriptionActive:   toNullBool(receiver.SubscriptionActive),
	})
}

func (r *SQLiteRepository) DeleteReceiver(ctx context.Context, id string) error {
	return r.queries.DeleteReceiver(ctx, id)
}

func (r *SQLiteRepository) GetReceiver(ctx context.Context, id string) (registry.Receiver, error) {
	rec, err := r.queries.GetReceiver(ctx, id)
	if err != nil {
		return registry.Receiver{}, err
	}
	return registry.Receiver{
		ID:                   rec.ID,
		ApiVersion:           rec.ApiVersion,
		DeviceID:             rec.DeviceID,
		Version:              rec.Version,
		Label:                rec.Label,
		Description:          rec.Description,
		Tags:                 rec.Tags,
		Transport:            rec.Transport,
		Format:               rec.Format,
		Caps:                 rec.Caps,
		InterfaceBindings:    rec.InterfaceBindings,
		SubscriptionSenderID: fromNullString(rec.SubscriptionSenderID),
		SubscriptionActive:   fromNullBool(rec.SubscriptionActive),
	}, nil
}

func (r *SQLiteRepository) ListReceivers(ctx context.Context) ([]registry.Receiver, error) {
	dbReceivers, err := r.queries.ListReceivers(ctx)
	if err != nil {
		return nil, err
	}
	receivers := make([]registry.Receiver, len(dbReceivers))
	for i, rec := range dbReceivers {
		receivers[i] = registry.Receiver{
			ID:                   rec.ID,
			ApiVersion:           rec.ApiVersion,
			DeviceID:             rec.DeviceID,
			Version:              rec.Version,
			Label:                rec.Label,
			Description:          rec.Description,
			Tags:                 rec.Tags,
			Transport:            rec.Transport,
			Format:               rec.Format,
			Caps:                 rec.Caps,
			InterfaceBindings:    rec.InterfaceBindings,
			SubscriptionSenderID: fromNullString(rec.SubscriptionSenderID),
			SubscriptionActive:   fromNullBool(rec.SubscriptionActive),
		}
	}
	return receivers, nil
}

func (r *SQLiteRepository) UpsertSubscription(ctx context.Context, subscription registry.Subscription) error {
	return r.queries.UpsertSubscription(ctx, db.UpsertSubscriptionParams{
		ID:              subscription.ID,
		ResourcePath:    subscription.ResourcePath,
		Params:          subscription.Params,
		MaxUpdateRateMs: toNullInt64(subscription.MaxUpdateRateMs),
		Persist:         toNullBool(subscription.Persist),
		SecureWebsocket: toNullBool(subscription.SecureWebsocket),
		WsHref:          toNullString(subscription.WsHref),
	})
}

func (r *SQLiteRepository) GetSubscription(ctx context.Context, id string) (registry.Subscription, error) {
	sub, err := r.queries.GetSubscription(ctx, id)
	if err != nil {
		return registry.Subscription{}, err
	}
	return registry.Subscription{
		ID:              sub.ID,
		ResourcePath:    sub.ResourcePath,
		Params:          sub.Params,
		MaxUpdateRateMs: fromNullInt64ToInt(sub.MaxUpdateRateMs),
		Persist:         fromNullBool(sub.Persist),
		SecureWebsocket: fromNullBool(sub.SecureWebsocket),
		WsHref:          fromNullString(sub.WsHref),
	}, nil
}

func (r *SQLiteRepository) ListSubscriptions(ctx context.Context) ([]registry.Subscription, error) {
	dbSubs, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	subscriptions := make([]registry.Subscription, len(dbSubs))
	for i, sub := range dbSubs {
		subscriptions[i] = registry.Subscription{
			ID:              sub.ID,
			ResourcePath:    sub.ResourcePath,
			Params:          sub.Params,
			MaxUpdateRateMs: fromNullInt64ToInt(sub.MaxUpdateRateMs),
			Persist:         fromNullBool(sub.Persist),
			SecureWebsocket: fromNullBool(sub.SecureWebsocket),
			WsHref:          fromNullString(sub.WsHref),
		}
	}
	return subscriptions, nil
}

func (r *SQLiteRepository) DeleteSubscription(ctx context.Context, id string) error {
	return r.queries.DeleteSubscription(ctx, id)
}

// Helpers
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nilToEmptyJSON(b []byte) []byte {
	if b == nil {
		return []byte("{}")
	}
	return b
}

func fromNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func toNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func fromNullInt64ToInt(ni sql.NullInt64) *int {
	if !ni.Valid {
		return nil
	}
	i := int(ni.Int64)
	return &i
}

func toNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func fromNullBool(nb sql.NullBool) *bool {
	if !nb.Valid {
		return nil
	}
	return &nb.Bool
}
