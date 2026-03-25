# Spotify API

Tiny API for Spotify "now playing" with OAuth login and callback.

## Run

```sh
cd /Users/harrybreen/Desktop/homelab/apps/spotify
SPOTIFY_CLIENT_ID=your-client-id \
SPOTIFY_SECRET=your-client-secret \
SPOTIFY_REDIRECT_URL=http://localhost:8080/spotify/auth/callback \
SPOTIFY_API_ADDR=:8080 \
go run .
```

## Authenticate

1. Make sure your Spotify app has a redirect URI that matches `SPOTIFY_REDIRECT_URL` (for example, `http://localhost:8080/spotify/auth/callback`).
2. Start the server with `SPOTIFY_CLIENT_ID`, `SPOTIFY_SECRET`, and `SPOTIFY_REDIRECT_URL` set.
3. Visit the login endpoint in your browser to begin OAuth:

```sh
open http://localhost:8080/spotify/auth/login
```

4. Approve the Spotify consent screen. Spotify will redirect you back to `/spotify/auth/callback` and the server will set a `spotify_session` cookie.
5. Use the cookie for API calls:

```sh
curl --cookie-jar /tmp/spotify.cookies http://localhost:8080/spotify/auth/login
curl --cookie /tmp/spotify.cookies http://localhost:8080/spotify/api/now-playing
```

## Token refresh

Access tokens are refreshed automatically when they are close to expiring. This requires Spotify to return a `refresh_token` during the initial OAuth login.

## Response fields

`/api/now-playing` returns:

- `is_playing`
- `track`
- `artists`
- `album_cover_url`
- `progress_ms`
- `duration_ms`

## Environment

- `SPOTIFY_CLIENT_ID`: Spotify client id (required).
- `SPOTIFY_SECRET`: Spotify client secret (required).
- `SPOTIFY_REDIRECT_URL`: OAuth callback URL (required).
- `SPOTIFY_SESSION_SECRET`: secret used to sign session cookies (required, 32+ characters).
- `SPOTIFY_API_ADDR`: listen address (default `:8080`).

## Kubernetes deploy

Set the image and secrets in `apps/spotify/manifests`:

- `apps/spotify/manifests/deployment.yaml`: update the image to your registry.
- `apps/spotify/manifests/secret.yaml`: fill in `SPOTIFY_CLIENT_ID`, `SPOTIFY_SECRET`, `SPOTIFY_REDIRECT_URL`, `SPOTIFY_SESSION_SECRET`, and `FRONTEND_REDIRECT_URL`.

Apply with kustomize:

```sh
kubectl apply -k /Users/harrybreen/Desktop/homelab/apps/spotify/manifests
```

If you use Argo CD, the application manifest is at `platform/argocd/applications/spotify.yaml`.
