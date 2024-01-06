package Project;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class PlaneRANSAC {

    //Declaring Instance Variables
    private PointCloud cloud;
    private double eps;
    private PointCloud domCloud;
    public Plane3D domPlane;
    public List<Point3D> domPoints;
    int bestSupp;


    //Declaring Getters and Setters
    public void setEPS(double eps) { this.eps = eps; }
    public double getEPS() { return eps; }

    //Constructor
    public PlaneRANSAC(PointCloud points) { this.cloud = points; }

    //Calculates the number of iterations for a c and p value.
    public int getNumberOfIterations(double c, double p) {
        
        return (int) (Math.log(1 - c) / Math.log(1 - Math.pow(p, 3)));
    }

    //Runs the RANSAC algorithm
    public void run (int numIterations, String file) {

        bestSupp = 0;
        domCloud = new PointCloud();

        for(int i = 0; i < numIterations; i++){

            //Creates a sample random plane from the original PointCloud
            List<Point3D> sample = new ArrayList<>();

            for(int j = 0; j < 3; j++){
                if(cloud.size() == 0){
                    System.out.println("Cloud is empty.");
                    return;
                }

                sample.add(cloud.getPoint());
            }

            Plane3D samplePlane = new Plane3D(sample.get(0), sample.get(1), sample.get(2));
            Iterator<Point3D> iter = cloud.iterator();

            int currSupp = 0;

            //Iterates through each point in the original PointCloud and counts the sample plane's supports
            while(iter.hasNext()){

                if(samplePlane.distanceFromPoint(iter.next()) < eps){

                    currSupp++;
                }
            }

            //Sets the dominant plane
            if(currSupp > bestSupp){

                domPoints = sample;
                bestSupp = currSupp;
                domPlane = samplePlane;
            }            
        }
        
        //Removes the dominant plane from the original PointCloud
        Iterator<Point3D> iter = cloud.iterator();
        while(iter.hasNext()){

            Point3D tmpPnt = iter.next();

            if(domPlane.distanceFromPoint(tmpPnt) < eps){

                iter.remove();
            }
        }

        for(Point3D p: domPoints){

            domCloud.addPoint(p);
        }


        domCloud.save(file);
        System.out.println("Done!");
    }
}