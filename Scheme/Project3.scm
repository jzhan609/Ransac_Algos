#lang scheme
;;Name: Jacob Zhang
;;Student ID: 300231094
;;Date: 4/12/2023
;;Purpose: Implement a simple RANSAC algorithm in scheme

; Outputs a Point Cloud represented by a list of Points given a fileName.
(define (readXYZ fileIn)
        (let ((sL (map (lambda s (string-split (car s)))
        (cdr (file->lines fileIn)))))
        (map (lambda (L)
        (map (lambda (s)
        (if (eqv? (string->number s) #f)
         s
(string->number s))) L)) sL)))

; Given 3 points, outputs the corresponding plane as a list of a, b, c, d.
(define (plane p1 p2 p3)
        (let* ((a (- (list-ref p1 1) (* (/ (- (list-ref p2 1) (list-ref p1 1)) (- (list-ref p2 0) (list-ref p1 0))) (list-ref p1 0))))
               (b (- (list-ref p1 2) (* (/ (- (list-ref p2 2) (list-ref p1 2)) (- (list-ref p2 0) (list-ref p1 0))) (list-ref p1 0))))
               (c (- (list-ref p1 2) (* (/ (- (list-ref p3 2) (list-ref p1 2)) (- (list-ref p3 0) (list-ref p1 0))) (list-ref p1 0))))
               (d (- (* a (list-ref p1 0)) (* b (list-ref p1 1)) (* c (list-ref p1 2)))))
        (list a b c d)))

; Given a Plane, Point Cloud and epsilon, outputs a pair consisting of the input plane and its number of supports.
(define (getSupport plane points eps)
        (let ((supportingPoints (map (lambda (p)
                                (distance p plane eps))
                                 points)))
        (cons (count (lambda (d) d) supportingPoints) plane)))

; Given confidence and percentage, calculates the number of iterations required k.
(define (getNumIterations c p)
        (ceiling (/ (log (- 1 c)) (log (- 1 (expt p 3))))))

; Given a point, plane and epsilon, outputs whether the point P is a support or not.
(define (distance p plane eps)
        (let* ((denominator (sqrt (+ (* (list-ref plane 0) (list-ref plane 0)) (* (list-ref plane 1) (list-ref plane 1)) (* (list-ref plane 2) (list-ref plane 2))))))
               (if (= denominator 0)
                    #f
               (let ((distance (/ (abs (+ (* (list-ref plane 0) (list-ref p 0)) (* (list-ref plane 1) (list-ref p 1)) (* (list-ref plane 2) (list-ref p 2)) (list-ref plane 3)))
                      denominator)))
               (> distance eps)))))

; Given a Point Cloud, k and epsilon, outputs the dominant plane of the point cloud.
(define (domPlane Pc k eps)
        (let ((maxPlane (getSupport (plane (list-ref Pc (random (length Pc)))
                                           (list-ref Pc (random (length Pc)))
                                           (list-ref Pc (random (length Pc))))
                                     Pc eps)))
        (let while ((k k) (currentMax maxPlane))
                   (if (= k 0)
                       (display currentMax)
                       (let ((current (getSupport (plane (list-ref Pc (random (length Pc)))
                                                         (list-ref Pc (random (length Pc)))
                                                         (list-ref Pc (random (length Pc))))
                                                   Pc eps)))
                   (while (- k 1)
                          (if (> (car current) (car currentMax))
                               current
                               currentMax)))))))

; The call function. Runs the program.
(define (planeRANSAC file c p eps)
        (let* ((Pc (readXYZ file))
               (k (getNumIterations c p)))
        (domPlane Pc k eps)))

;Run Commands, Uncomment to use, File must be in same directory.
(planeRANSAC "Point_Cloud_1_No_Road_Reduced.xyz" 0.95 0.099 0.099)
;(planeRANSAC "Point_Cloud_2_No_Road_Reduced.xyz" 0.95 0.099 0.099)
;(planeRANSAC "Point_Cloud_3_No_Road_Reduced.xyz" 0.95 0.099 0.099)