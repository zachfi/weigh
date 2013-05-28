module Weigh
  module Report

    def self.print(summary)
      summary.sort{|a,b| a[1]<=>b[1]}.each { |elem|
        size     = elem[1]
        filename = elem[0]
        puts sprintf("%15s   %s\n", neat_size(size), filename)
      }

      puts sprintf("%16s %s\n", "---", "---")
      puts sprintf("%15s   %s\n", neat_size(total_size), ":total size")
      puts sprintf("%16s %s\n", "---", "---")
    end
  end
end
