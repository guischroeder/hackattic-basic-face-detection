# Hackattic Basic Face Detection Challenge

A Clojure solution for the [Hackattic Basic Face Detection Challenge](https://hackattic.com/challenges/basic_face_detection).

## Challenge Description

The challenge requires detecting faces in an image provided by the Hackattic API and submitting the coordinates of the detected faces. The solution uses:

- Clojure for the programming language
- OpenCV (via JavaCV) for face detection
- Haar cascade classifiers for the face detection algorithm

## Prerequisites

- [Java JDK](https://adoptopenjdk.net/) (version 8 or higher)
- [Leiningen](https://leiningen.org/) (Clojure build tool)
- A Hackattic account and access token

## Dependencies

This project depends on:
- clj-http for HTTP requests
- cheshire for JSON parsing
- JavaCV/OpenCV for face detection

All dependencies are managed by Leiningen and will be automatically downloaded when building the project.

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/guischroeder/hackattic-basic-face-detection.git
   cd hackattic-basic-face-detection
   ```

2. Install dependencies:
   ```bash
   lein deps
   ```

## Getting a Hackattic Access Token

1. Sign up/login to [Hackattic](https://hackattic.com/)
2. Navigate to the [Basic Face Detection Challenge](https://hackattic.com/challenges/basic_face_detection)
3. Your access token is in the URL or can be found in your profile

## Running the Solution

Run the solution with your Hackattic access token:

```bash
lein run YOUR_ACCESS_TOKEN
```

Or set it as an environment variable:

```bash
export HACKATTIC_TOKEN=your_token_here
lein run
```

The program will:
1. Fetch the problem from Hackattic API
2. Download the image containing faces
3. Detect faces using OpenCV
4. Submit the coordinates of detected faces back to Hackattic
5. Display the result

## Project Structure

- `src/hackattic_basic_face_detection/`
  - `api/hackattic.clj` - API client for Hackattic
  - `detection/face.clj` - Face detection implementation
  - `challenge.clj` - Challenge workflow orchestration
  - `core.clj` - Main entry point
  - `util/` - Utility namespaces for OpenCV
- `resources/` - Resource files including downloaded cascade files

## How It Works

1. The application fetches a problem from Hackattic containing an image URL
2. It downloads the image to the resources directory
3. OpenCV with Haar cascade classifiers is used to detect faces in the image
4. The coordinates of the detected faces are collected in the format expected by Hackattic
5. The solution is submitted back to Hackattic for verification

## Troubleshooting

### OpenCV Issues

If you encounter OpenCV-related errors:

1. Ensure your Java version is compatible (OpenCV works best with Java 8-11)

2. Check that the JavaCV dependencies were downloaded correctly:
   ```bash
   lein deps :tree
   ```

3. If problems persist, try running with expanded memory:
   ```bash
   JAVA_OPTS="-Xmx1g" lein run YOUR_ACCESS_TOKEN
   ```

### API Issues

If you encounter API-related errors:

1. Verify your access token is correct
2. Ensure you have an active internet connection
3. Check that the Hackattic API is available

## License

Copyright © 2025

This program and the accompanying materials are made available under the terms of the Eclipse Public License 2.0.
