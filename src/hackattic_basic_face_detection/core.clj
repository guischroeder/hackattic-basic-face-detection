(ns hackattic-basic-face-detection.core
  (:require [hackattic-basic-face-detection.hackattic.client :as client]
            [hackattic-basic-face-detection.face_detection.face_detector :as face])
  (:gen-class))

(defn get-access-token
  "Gets the access token from args or environment"
  [args]
  (or (first args)
      (System/getenv "HACKATTIC_TOKEN")
      (throw (IllegalArgumentException.
              "Access token required: provide as arg or HACKATTIC_TOKEN env var"))))

(defn solve-face-detection
  "Orchestrates the workflow for solving the face detection challenge"
  [access-token]
  (println "Fetching problem...")
  (let [problem (client/get-problem access-token)
        image-url (get problem :image_url)

        _ (println "Downloading image from" image-url)
        image-path (client/download-image image-url "resources/challenge-image.jpg")

        _ (println "Detecting faces...")
        face-coordinates (face/detect-faces image-path)

        _ (println "Submitting solution with" (count face-coordinates) "faces...")
        result (client/submit-solution access-token face-coordinates)]

    (println "Result from Hackattic:" result)
    result))

(defn -main
  "Entry point for the Hackattic face detection challenge solution"
  [& args]
  (try
    (let [access-token (get-access-token args)]
      (solve-face-detection access-token)
      (System/exit 0))
    (catch Exception e
      (println "Error:" (.getMessage e))
      (System/exit 1))))
