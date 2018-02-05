CREATE TABLE links
(
  CONSTRAINT links_pk PRIMARY KEY (id),
  id INT8 NOT NULL DEFAULT generate_primary_key(),
  url TEXT NOT NULL,
  token TEXT NOT NULL,
  conversion INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc') NOT NULL
);
CREATE UNIQUE INDEX links_url_uindex ON links (url);
CREATE UNIQUE INDEX links_token_uindex ON links (token);
