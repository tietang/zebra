package zk.mock;

import org.apache.curator.test.TestingServer;

import java.io.IOException;

/**
 * @author Tietang Wang 铁汤
 * @Date: 2017/8/5 下午11:15
 */
public class Server {

    public static void main(String[] args) throws Exception {
        int port = 2181;
        if (args.length >= 1) {
            port = Integer.parseInt(args[0]);
        }
        final TestingServer zkServer = new TestingServer(2181, true);

        zkServer.start();
        System.out.println("started");
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                try {
                    zkServer.stop();
                } catch (IOException e) {
                    e.printStackTrace();

                }
            }
        });
    }
}
