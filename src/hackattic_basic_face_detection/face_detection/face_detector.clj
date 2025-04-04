(ns hackattic-basic-face-detection.face-detection.face-detector
  (:require [clojure.java.io :as io])
  (:import [org.bytedeco.opencv.opencv_core Mat Rect RectVector Size]
           [org.bytedeco.opencv.opencv_imgcodecs Imgcodecs]
           [org.bytedeco.opencv.opencv_imgproc Imgproc]
           [org.bytedeco.opencv.opencv_objdetect CascadeClassifier]
           [java.io File]))

(def cascade-file-path
  (or (System/getenv "CASCADE_FILE_PATH")
      "resources/haarcascade_frontalface_default.xml"))

(def detection-params
  {:scale-factor 1.1    ; How much the image size is reduced at each image scale
   :min-neighbors 3     ; How many neighbors each candidate rectangle should have
   :min-face-size 30    ; Minimum possible face size in pixels
   :flags 0})           ; Not used in newer OpenCV versions

(defn ensure-cascade-file
  "Ensures the Haar cascade file exists, downloading if necessary"
  []
  (let [cascade-file (File. cascade-file-path)]
    (when-not (.exists cascade-file)
      (println "Downloading Haar cascade file...")
      (io/make-parents cascade-file-path)
      (let [url "https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalface_default.xml"]
        (with-open [in (io/input-stream url)
                    out (io/output-stream cascade-file-path)]
          (io/copy in out))
        (println "Haar cascade file downloaded to" cascade-file-path)))))

(defn load-and-preprocess-image
  "Loads and preprocesses an image for better face detection"
  [image-path]
  (let [image (Imgcodecs/imread image-path)
        gray-image (Mat.)]
    ; Convert to grayscale (improves detection)
    (Imgproc/cvtColor image gray-image Imgproc/COLOR_BGR2GRAY)
    ; Equalize histogram (helps with varying lighting)
    (Imgproc/equalizeHist gray-image gray-image)
    gray-image))

(defn detect-faces-in-image
  "Performs the actual face detection on the preprocessed image"
  [gray-image]
  (let [{:keys [scale-factor min-neighbors flags min-face-size]} detection-params
        face-detector (CascadeClassifier. cascade-file-path)
        faces (RectVector.)]
    (.detectMultiScale face-detector
                       gray-image
                       faces
                       scale-factor
                       min-neighbors
                       flags
                       (Size. min-face-size min-face-size)
                       (Size.))
    faces))

(defn convert-to-coordinates
  "Converts OpenCV Rect objects to maps with coordinates"
  [faces]
  (mapv (fn [^Rect rect]
          {:x (.x rect)
           :y (.y rect)
           :width (.width rect)
           :height (.height rect)})
        (iterator-seq (.iterator faces))))

(defn detect-faces
  "Main function that detects faces in the image and returns their coordinates"
  [image-path]
  (try
    (ensure-cascade-file)

    (let [preprocessed-image (load-and-preprocess-image image-path)
          detected-faces (detect-faces-in-image preprocessed-image)
          face-coordinates (convert-to-coordinates detected-faces)]

      (if (empty? face-coordinates)
        (println "Warning: No faces detected in the image!")
        (println "Detected" (count face-coordinates) "faces"))

      face-coordinates)

    (catch Exception e
      (println "Error in face detection:" (.getMessage e))
      (.printStackTrace e)
      (System/exit 1))))
