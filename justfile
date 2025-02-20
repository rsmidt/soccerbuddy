generate-proto:
    @echo "Generating proto files..."
    buf generate
    cp -r gen/web/. web/src/lib/gen/
    cp -r gen/web/. mobile/api/

generate-go:
    @echo "Generating using go generate..."
    go generate ./...

generate: generate-proto generate-go

copy-permify-schema:
    cat permify/schema.perm | jq -Rsa | wl-copy

[working-directory: 'mobile']
build-android:
    GOOGLE_SERVICES_JSON="{{justfile_directory()}}/mobile/google-services.json" \
    GOOGLE_SERVICES_PLIST="{{justfile_directory()}}/mobile/GoogleService-Info.plist" \
    EXPO_PUBLIC_API_URL="https://soccerbuddy.app/api" \
    EXPO_PUBLIC_URL="https://soccerbuddy.app" \
    eas build --platform android --profile production --local

[working-directory: 'mobile']
build-ios:
    GOOGLE_SERVICES_JSON="{{justfile_directory()}}/mobile/google-services.json" \
    GOOGLE_SERVICES_PLIST="{{justfile_directory()}}/mobile/GoogleService-Info.plist" \
    EXPO_PUBLIC_API_URL="https://soccerbuddy.app/api" \
    EXPO_PUBLIC_URL="https://soccerbuddy.app" \
    eas build --platform ios --profile production --local
