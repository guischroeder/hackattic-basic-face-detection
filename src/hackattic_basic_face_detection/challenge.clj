(ns hackattic-basic-face-detection.challenge
  (:require [hackattic-basic-face-detection.client.hackattic :as client]
            [hackattic-basic-face-detection.face_detection.face_detector :as face]))

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
