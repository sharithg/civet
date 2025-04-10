#!/bin/bash

set -e

source ./.secret/env.sh

# === CONFIGURATION ===
OWNER="sharithg"
REPO="civet"
IMAGE_NAME="civet"

# Generate a truly random string (e.g., 8 alphanumeric chars)
RANDOM_TAG=$(LC_ALL=C tr -dc a-z0-9 </dev/urandom | head -c 8)
TAG="local-${RANDOM_TAG}"
GHCR_IMAGE="ghcr.io/${OWNER}/${REPO}:${TAG}"

# Switch to default Docker context for build/push
docker context use default

echo "🔨 Building Docker image with tag: ${TAG}..."
docker build --platform linux/amd64 -t ${IMAGE_NAME}:${TAG} .

echo "🏷️ Tagging image as ${GHCR_IMAGE}..."
docker tag ${IMAGE_NAME}:${TAG} ${GHCR_IMAGE}

echo "🔐 Logging into GitHub Container Registry..."
echo "${GITHUB_TOKEN}" | docker login ghcr.io -u ${OWNER} --password-stdin

echo "📤 Pushing image to GitHub Container Registry..."
docker push ${GHCR_IMAGE}

echo "✅ Done! Image pushed as ${GHCR_IMAGE}"

export GIT_COMMIT_HASH=${TAG}

# Switch to deploy context and deploy
docker context use civet-app
docker stack deploy -c docker-stack.yaml civet

echo "🚀 Deployment complete using image tag: ${TAG}"
