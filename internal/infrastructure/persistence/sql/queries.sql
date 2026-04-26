-- name: CreateNode :exec
INSERT INTO nodes (
    id, api_version, version, label, description, tags, href, hostname, api, caps, services, clocks, interfaces, last_seen
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description), sqlc.arg(tags),
    sqlc.arg(href), sqlc.arg(hostname), sqlc.arg(api), sqlc.arg(caps), sqlc.arg(services),
    sqlc.arg(clocks), sqlc.arg(interfaces), sqlc.arg(last_seen)
);

-- name: UpsertNode :exec
INSERT INTO nodes (
    id, api_version, version, label, description, tags, href, hostname, api, caps, services, clocks, interfaces, last_seen
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description), sqlc.arg(tags),
    sqlc.arg(href), sqlc.arg(hostname), sqlc.arg(api), sqlc.arg(caps), sqlc.arg(services),
    sqlc.arg(clocks), sqlc.arg(interfaces), sqlc.arg(last_seen)
)
ON CONFLICT(id) DO UPDATE SET
    api_version = excluded.api_version, version = excluded.version, label = excluded.label, description = excluded.description,
    tags = excluded.tags, href = excluded.href, hostname = excluded.hostname,
    api = excluded.api, caps = excluded.caps, services = excluded.services,
    clocks = excluded.clocks, interfaces = excluded.interfaces, last_seen = excluded.last_seen;

-- name: GetNode :one
SELECT * FROM nodes WHERE id = sqlc.arg(id);

-- name: ListNodes :many
SELECT * FROM nodes;

-- name: UpdateNode :exec
UPDATE nodes SET
    version = sqlc.arg(version), label = sqlc.arg(label), description = sqlc.arg(description),
    tags = sqlc.arg(tags), href = sqlc.arg(href), hostname = sqlc.arg(hostname),
    api = sqlc.arg(api), caps = sqlc.arg(caps), services = sqlc.arg(services),
    clocks = sqlc.arg(clocks), interfaces = sqlc.arg(interfaces), last_seen = sqlc.arg(last_seen)
WHERE id = sqlc.arg(id);

-- name: DeleteNode :exec
DELETE FROM nodes WHERE id = sqlc.arg(id);

-- name: UpdateNodeLastSeen :exec
UPDATE nodes SET last_seen = sqlc.arg(last_seen) WHERE id = sqlc.arg(id);

-- name: GetExpiredNodes :many
SELECT id FROM nodes WHERE last_seen < sqlc.arg(expiration_time);

-- name: CreateDevice :exec
INSERT INTO devices (
    id, api_version, node_id, version, label, description, tags, type, senders, receivers, controls
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(node_id), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description),
    sqlc.arg(tags), sqlc.arg(type), sqlc.arg(senders), sqlc.arg(receivers), sqlc.arg(controls)
);

-- name: UpsertDevice :exec
INSERT INTO devices (
    id, api_version, node_id, version, label, description, tags, type, senders, receivers, controls
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(node_id), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description),
    sqlc.arg(tags), sqlc.arg(type), sqlc.arg(senders), sqlc.arg(receivers), sqlc.arg(controls)
)
ON CONFLICT(id) DO UPDATE SET
    api_version = excluded.api_version, node_id = excluded.node_id, version = excluded.version, label = excluded.label,
    description = excluded.description, tags = excluded.tags, type = excluded.type,
    senders = excluded.senders, receivers = excluded.receivers, controls = excluded.controls;

-- name: GetDevice :one
SELECT * FROM devices WHERE id = sqlc.arg(id);

-- name: ListDevices :many
SELECT * FROM devices;

-- name: ListDevicesByNode :many
SELECT * FROM devices WHERE node_id = sqlc.arg(node_id);

-- name: UpdateDevice :exec
UPDATE devices SET
    version = sqlc.arg(version), label = sqlc.arg(label), description = sqlc.arg(description),
    tags = sqlc.arg(tags), type = sqlc.arg(type), senders = sqlc.arg(senders),
    receivers = sqlc.arg(receivers), controls = sqlc.arg(controls)
WHERE id = sqlc.arg(id);

-- name: DeleteDevice :exec
DELETE FROM devices WHERE id = sqlc.arg(id);

-- name: CreateSource :exec
INSERT INTO sources (
    id, api_version, device_id, version, label, description, tags, grain_rate, format, caps, parents, clock_name
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(device_id), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description),
    sqlc.arg(tags), sqlc.arg(grain_rate), sqlc.arg(format), sqlc.arg(caps), sqlc.arg(parents), sqlc.arg(clock_name)
);

-- name: UpsertSource :exec
INSERT INTO sources (
    id, api_version, device_id, version, label, description, tags, grain_rate, format, caps, parents, clock_name
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(device_id), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description),
    sqlc.arg(tags), sqlc.arg(grain_rate), sqlc.arg(format), sqlc.arg(caps), sqlc.arg(parents), sqlc.arg(clock_name)
)
ON CONFLICT(id) DO UPDATE SET
    api_version = excluded.api_version, device_id = excluded.device_id, version = excluded.version, label = excluded.label,
    description = excluded.description, tags = excluded.tags, grain_rate = excluded.grain_rate,
    format = excluded.format, caps = excluded.caps, parents = excluded.parents,
    clock_name = excluded.clock_name;

-- name: GetSource :one
SELECT * FROM sources WHERE id = sqlc.arg(id);

-- name: ListSources :many
SELECT * FROM sources;

-- name: ListSourcesByDevice :many
SELECT * FROM sources WHERE device_id = sqlc.arg(device_id);

-- name: UpdateSource :exec
UPDATE sources SET
    version = sqlc.arg(version), label = sqlc.arg(label), description = sqlc.arg(description),
    tags = sqlc.arg(tags), grain_rate = sqlc.arg(grain_rate), format = sqlc.arg(format),
    caps = sqlc.arg(caps), parents = sqlc.arg(parents), clock_name = sqlc.arg(clock_name)
WHERE id = sqlc.arg(id);

-- name: DeleteSource :exec
DELETE FROM sources WHERE id = sqlc.arg(id);

-- name: CreateFlow :exec
INSERT INTO flows (
    id, api_version, source_id, device_id, version, label, description, tags, format, media_type,
    sample_rate, bit_depth, DID_SDID, grain_rate, frame_width, frame_height,
    interlace_mode, colorspace, components, transfer_characteristic
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(source_id), sqlc.arg(device_id), sqlc.arg(version), sqlc.arg(label),
    sqlc.arg(description), sqlc.arg(tags), sqlc.arg(format), sqlc.arg(media_type),
    sqlc.arg(sample_rate), sqlc.arg(bit_depth), sqlc.arg(DID_SDID), sqlc.arg(grain_rate),
    sqlc.arg(frame_width), sqlc.arg(frame_height), sqlc.arg(interlace_mode), sqlc.arg(colorspace),
    sqlc.arg(components), sqlc.arg(transfer_characteristic)
);

-- name: UpsertFlow :exec
INSERT INTO flows (
    id, api_version, source_id, device_id, version, label, description, tags, format, media_type,
    sample_rate, bit_depth, DID_SDID, grain_rate, frame_width, frame_height,
    interlace_mode, colorspace, components, transfer_characteristic
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(source_id), sqlc.arg(device_id), sqlc.arg(version), sqlc.arg(label),
    sqlc.arg(description), sqlc.arg(tags), sqlc.arg(format), sqlc.arg(media_type),
    sqlc.arg(sample_rate), sqlc.arg(bit_depth), sqlc.arg(DID_SDID), sqlc.arg(grain_rate),
    sqlc.arg(frame_width), sqlc.arg(frame_height), sqlc.arg(interlace_mode), sqlc.arg(colorspace),
    sqlc.arg(components), sqlc.arg(transfer_characteristic)
)
ON CONFLICT(id) DO UPDATE SET
    api_version = excluded.api_version, source_id = excluded.source_id, device_id = excluded.device_id, version = excluded.version,
    label = excluded.label, description = excluded.description, tags = excluded.tags,
    format = excluded.format, media_type = excluded.media_type, sample_rate = excluded.sample_rate,
    bit_depth = excluded.bit_depth, DID_SDID = excluded.DID_SDID, grain_rate = excluded.grain_rate,
    frame_width = excluded.frame_width, frame_height = excluded.frame_height,
    interlace_mode = excluded.interlace_mode, colorspace = excluded.colorspace,
    components = excluded.components, transfer_characteristic = excluded.transfer_characteristic;

-- name: GetFlow :one
SELECT * FROM flows WHERE id = sqlc.arg(id);

-- name: ListFlows :many
SELECT * FROM flows;

-- name: ListFlowsByDevice :many
SELECT * FROM flows WHERE device_id = sqlc.arg(device_id);

-- name: ListFlowsBySource :many
SELECT * FROM flows WHERE source_id = sqlc.arg(source_id);

-- name: UpdateFlow :exec
UPDATE flows SET
    version = sqlc.arg(version), label = sqlc.arg(label), description = sqlc.arg(description),
    tags = sqlc.arg(tags), format = sqlc.arg(format), media_type = sqlc.arg(media_type),
    sample_rate = sqlc.arg(sample_rate), bit_depth = sqlc.arg(bit_depth),
    DID_SDID = sqlc.arg(DID_SDID), grain_rate = sqlc.arg(grain_rate),
    frame_width = sqlc.arg(frame_width), frame_height = sqlc.arg(frame_height),
    interlace_mode = sqlc.arg(interlace_mode), colorspace = sqlc.arg(colorspace),
    components = sqlc.arg(components), transfer_characteristic = sqlc.arg(transfer_characteristic)
WHERE id = sqlc.arg(id);

-- name: DeleteFlow :exec
DELETE FROM flows WHERE id = sqlc.arg(id);

-- name: CreateSender :exec
INSERT INTO senders (
    id, api_version, device_id, flow_id, version, label, description, tags, transport,
    manifest_href, interface_bindings, subscription_receiver_id, subscription_active
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(device_id), sqlc.arg(flow_id), sqlc.arg(version), sqlc.arg(label),
    sqlc.arg(description), sqlc.arg(tags), sqlc.arg(transport), sqlc.arg(manifest_href),
    sqlc.arg(interface_bindings), sqlc.arg(subscription_receiver_id), sqlc.arg(subscription_active)
);

-- name: UpsertSender :exec
INSERT INTO senders (
    id, api_version, device_id, flow_id, version, label, description, tags, transport,
    manifest_href, interface_bindings, subscription_receiver_id, subscription_active
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(device_id), sqlc.arg(flow_id), sqlc.arg(version), sqlc.arg(label),
    sqlc.arg(description), sqlc.arg(tags), sqlc.arg(transport), sqlc.arg(manifest_href),
    sqlc.arg(interface_bindings), sqlc.arg(subscription_receiver_id), sqlc.arg(subscription_active)
)
ON CONFLICT(id) DO UPDATE SET
    api_version = excluded.api_version, device_id = excluded.device_id, flow_id = excluded.flow_id, version = excluded.version,
    label = excluded.label, description = excluded.description, tags = excluded.tags,
    transport = excluded.transport, manifest_href = excluded.manifest_href,
    interface_bindings = excluded.interface_bindings,
    subscription_receiver_id = excluded.subscription_receiver_id,
    subscription_active = excluded.subscription_active;

-- name: GetSender :one
SELECT * FROM senders WHERE id = sqlc.arg(id);

-- name: ListSenders :many
SELECT * FROM senders;

-- name: ListSendersByDevice :many
SELECT * FROM senders WHERE device_id = sqlc.arg(device_id);

-- name: ListSendersByFlow :many
SELECT * FROM senders WHERE flow_id = sqlc.arg(flow_id);

-- name: UpdateSender :exec
UPDATE senders SET
    version = sqlc.arg(version), label = sqlc.arg(label), description = sqlc.arg(description),
    tags = sqlc.arg(tags), flow_id = sqlc.arg(flow_id), transport = sqlc.arg(transport),
    manifest_href = sqlc.arg(manifest_href), interface_bindings = sqlc.arg(interface_bindings),
    subscription_receiver_id = sqlc.arg(subscription_receiver_id),
    subscription_active = sqlc.arg(subscription_active)
WHERE id = sqlc.arg(id);

-- name: DeleteSender :exec
DELETE FROM senders WHERE id = sqlc.arg(id);

-- name: CreateReceiver :exec
INSERT INTO receivers (
    id, api_version, device_id, version, label, description, tags, transport, format, caps,
    interface_bindings, subscription_sender_id, subscription_active
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(device_id), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description),
    sqlc.arg(tags), sqlc.arg(transport), sqlc.arg(format), sqlc.arg(caps),
    sqlc.arg(interface_bindings), sqlc.arg(subscription_sender_id), sqlc.arg(subscription_active)
);

-- name: UpsertReceiver :exec
INSERT INTO receivers (
    id, api_version, device_id, version, label, description, tags, transport, format, caps,
    interface_bindings, subscription_sender_id, subscription_active
) VALUES (
    sqlc.arg(id), sqlc.arg(api_version), sqlc.arg(device_id), sqlc.arg(version), sqlc.arg(label), sqlc.arg(description),
    sqlc.arg(tags), sqlc.arg(transport), sqlc.arg(format), sqlc.arg(caps),
    sqlc.arg(interface_bindings), sqlc.arg(subscription_sender_id), sqlc.arg(subscription_active)
)
ON CONFLICT(id) DO UPDATE SET
    api_version = excluded.api_version, device_id = excluded.device_id, version = excluded.version, label = excluded.label,
    description = excluded.description, tags = excluded.tags, transport = excluded.transport,
    format = excluded.format, caps = excluded.caps, interface_bindings = excluded.interface_bindings,
    subscription_sender_id = excluded.subscription_sender_id,
    subscription_active = excluded.subscription_active;

-- name: GetReceiver :one
SELECT * FROM receivers WHERE id = sqlc.arg(id);

-- name: ListReceivers :many
SELECT * FROM receivers;

-- name: ListReceiversByDevice :many
SELECT * FROM receivers WHERE device_id = sqlc.arg(device_id);

-- name: UpdateReceiver :exec
UPDATE receivers SET
    version = sqlc.arg(version), label = sqlc.arg(label), description = sqlc.arg(description),
    tags = sqlc.arg(tags), transport = sqlc.arg(transport), format = sqlc.arg(format),
    caps = sqlc.arg(caps), interface_bindings = sqlc.arg(interface_bindings),
    subscription_sender_id = sqlc.arg(subscription_sender_id),
    subscription_active = sqlc.arg(subscription_active)
WHERE id = sqlc.arg(id);

-- name: DeleteReceiver :exec
DELETE FROM receivers WHERE id = sqlc.arg(id);

-- name: CreateSubscription :exec
INSERT INTO subscriptions (
    id, resource_path, params, max_update_rate_ms, persist, secure_websocket, ws_href
) VALUES (
    sqlc.arg(id), sqlc.arg(resource_path), sqlc.arg(params), sqlc.arg(max_update_rate_ms),
    sqlc.arg(persist), sqlc.arg(secure_websocket), sqlc.arg(ws_href)
);

-- name: UpsertSubscription :exec
INSERT INTO subscriptions (
    id, resource_path, params, max_update_rate_ms, persist, secure_websocket, ws_href
) VALUES (
    sqlc.arg(id), sqlc.arg(resource_path), sqlc.arg(params), sqlc.arg(max_update_rate_ms),
    sqlc.arg(persist), sqlc.arg(secure_websocket), sqlc.arg(ws_href)
)
ON CONFLICT(id) DO UPDATE SET
    resource_path = excluded.resource_path, params = excluded.params,
    max_update_rate_ms = excluded.max_update_rate_ms, persist = excluded.persist,
    secure_websocket = excluded.secure_websocket, ws_href = excluded.ws_href;

-- name: GetSubscription :one
SELECT * FROM subscriptions WHERE id = sqlc.arg(id);

-- name: ListSubscriptions :many
SELECT * FROM subscriptions;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE id = sqlc.arg(id);

-- name: IDExistsInNodes :one
SELECT COUNT(*) FROM nodes WHERE id = sqlc.arg(id);

-- name: IDExistsInDevices :one
SELECT COUNT(*) FROM devices WHERE id = sqlc.arg(id);

-- name: IDExistsInSources :one
SELECT COUNT(*) FROM sources WHERE id = sqlc.arg(id);

-- name: IDExistsInFlows :one
SELECT COUNT(*) FROM flows WHERE id = sqlc.arg(id);

-- name: IDExistsInSenders :one
SELECT COUNT(*) FROM senders WHERE id = sqlc.arg(id);

-- name: IDExistsInReceivers :one
SELECT COUNT(*) FROM receivers WHERE id = sqlc.arg(id);
