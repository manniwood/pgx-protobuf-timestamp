# 2.0.0 (2025-05-14)

* Upgrades pgx to v5.7.4
* Handles pgx v5.6.0 API change where `pgtype.Type.Codec` used to require (tolerate?) a `TimestampCodec{}` or `TimestampTZCoded{}` but now requires `&TimestampCodec{}` and `&TimestampTZCoded{}` (Thanks Dmitri Dolguikh / dmitri-d for reporting this change.)

# 1.2.0 (2025-05-14)

* Upgrades pgx to v5.5.5

# 1.1.1

* Fixes bp/pb transposition in docs

# 1.1.0

* Adds better docs

# 1.0.0

* Inaugural release
