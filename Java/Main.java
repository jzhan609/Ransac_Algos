package Project;

import java.io.IOException;
import java.util.Scanner;

public class Main {

    public static void main(String[] args) throws IOException {
        
        Scanner sc = new Scanner(System.in);
        System.out.println("File name: ");
        String fileString = sc.nextLine();
        PointCloud pc = new PointCloud(fileString);
        PlaneRANSAC pr = new PlaneRANSAC(pc);
        System.out.println("Value of c: ");
        double c = Double.parseDouble(sc.nextLine());
        System.out.println("Value of p: ");
        double p = Double.parseDouble(sc.nextLine());
        System.out.println("Epsilon value: ");
        pr.setEPS(Double.parseDouble(sc.nextLine()));

        int numIterations = pr.getNumberOfIterations(c, p);
        for(int i = 1; i < 4; i++){

            String outFileString = fileString.replace(".xyz", "_" + i + ".xyz");
            pr.run(numIterations, outFileString);
        }
        
        sc.close();
    }
}
