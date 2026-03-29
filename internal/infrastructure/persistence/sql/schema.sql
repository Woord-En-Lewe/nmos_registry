CREATE TABLE nodes (
    id TEXT PRIMARY KEY,
    version TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    tags JSON,
    href TEXT NOT NULL,
    hostname TEXT,
    api JSON NOT NULL,
    caps JSON,
    services JSON,
    clocks JSON,
    interfaces JSON,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    tags JSON,
    type TEXT NOT NULL,
    senders JSON,
    receivers JSON,
    controls JSON
);

CREATE TABLE sources (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    tags JSON,
    grain_rate JSON NOT NULL,
    format TEXT NOT NULL,
    caps JSON,
    parents JSON,
    clock_name TEXT
);

CREATE TABLE flows (
    id TEXT PRIMARY KEY,
    source_id TEXT NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    tags JSON,
    format TEXT NOT NULL,
    media_type TEXT,
    sample_rate JSON,
    bit_depth INTEGER,
    DID_SDID JSON,
    grain_rate JSON,
    frame_width INTEGER,
    frame_height INTEGER,
    interlace_mode TEXT,
    colorspace TEXT,
    components JSON,
    transfer_characteristic TEXT
);

CREATE TABLE senders (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    flow_id TEXT REFERENCES flows(id) ON DELETE SET NULL,
    version TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    tags JSON,
    transport TEXT NOT NULL,
    manifest_href TEXT,
    interface_bindings JSON,
    subscription_receiver_id TEXT,
    subscription_active BOOLEAN
);

CREATE TABLE receivers (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    tags JSON,
    transport TEXT NOT NULL,
    format TEXT NOT NULL,
    caps JSON,
    interface_bindings JSON,
    subscription_sender_id TEXT,
    subscription_active BOOLEAN
);

CREATE TABLE subscriptions (
    id TEXT PRIMARY KEY,
    resource_path TEXT NOT NULL,
    params JSON,
    max_update_rate_ms INTEGER,
    persist BOOLEAN,
    secure_websocket BOOLEAN,
    ws_href TEXT
);
