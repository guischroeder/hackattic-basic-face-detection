(ns hackattic-basic-face-detection.core
  (:require [hackattic-basic-face-detection.challenge :as challenge])
  (:gen-class))

(defn get-access-token
  "Gets the access token from args or environment"
  [args]
  (or (first args)
      (System/getenv "HACKATTIC_TOKEN")
      (throw (IllegalArgumentException.
              "Access token required: provide as arg or HACKATTIC_TOKEN env var"))))

(defn -main
  "Entry point for the Hackattic face detection challenge solution"
  [& args]
  (try
    (let [access-token (get-access-token args)]
      (challenge/solve-face-detection access-token)
      (System/exit 0))
    (catch Exception e
      (println "Error:" (.getMessage e))
      (System/exit 1))))
