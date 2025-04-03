(ns hackattic-basic-face-detection.api.hackattic
  (:require [clj-http.client :as http]
            [cheshire.core :as json]
            [clojure.java.io :as io]))

(def base-url "https://hackattic.com/challenges/basic_face_detection")

(defn- handle-api-error
  "Generic error handler for API calls"
  [response operation]
  (when (>= (:status response) 400)
    (throw (ex-info (str "Failed to " operation)
                    {:status (:status response)
                     :body (:body response)}))))

(defn- build-url
  "Builds a URL for a specific endpoint with access token"
  [endpoint access-token]
  (str base-url endpoint "?access_token=" access-token))

(defn get-problem
  "Fetches the face detection problem from Hackattic API"
  [access-token]
  (try
    (let [url (build-url "/problem" access-token)
          response (http/get url {:as :json
                                  :throw-exceptions false})]
      (handle-api-error response "fetch problem")
      (:body response))
    (catch Exception e
      (println "Error fetching problem:" (.getMessage e))
      (when (ex-data e)
        (println "Details:" (ex-data e)))
      (System/exit 1))))

(defn download-image
  "Downloads an image from URL to the specified path"
  [image-url output-path]
  (try
    (io/make-parents output-path)
    (with-open [in (io/input-stream image-url)
                out (io/output-stream output-path)]
      (io/copy in out))
    output-path
    (catch Exception e
      (println "Error downloading image:" (.getMessage e))
      (System/exit 1))))

(defn submit-solution
  "Submits the face coordinates solution to Hackattic API"
  [access-token face-coordinates]
  (try
    (let [url (build-url "/solve" access-token)
          solution {:faces face-coordinates}
          response (http/post url {:body (json/generate-string solution)
                                   :content-type :json
                                   :as :json
                                   :throw-exceptions false})]
      (handle-api-error response "submit solution")
      (:body response))
    (catch Exception e
      (println "Error submitting solution:" (.getMessage e))
      (when (ex-data e)
        (println "Details:" (ex-data e)))
      (System/exit 1))))
