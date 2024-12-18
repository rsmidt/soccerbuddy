generate-proto:
    buf generate
    cp -r gen/web/. web/src/lib/gen/
    cp -r gen/web/. mobile/api/

copy-permify-schema:
    cat permify/schema.perm | jq -Rsa | wl-copy
