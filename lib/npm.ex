defmodule Npm do

  @depFileName "package.json"
  def getDepFileName do
    @depFileName
  end

  def getDockerFile do
    """
    FROM kkarczmarczyk/node-yarn:8.0-wheezy
    
    WORKDIR /build/app
    
    ENV PATH=/build/node_modules/.bin:$PATH
    
    ADD package.json /build/
    
    RUN yarn && chmod -R 777 /build
    
    RUN mkdir /.config /.cache && chmod -R 777 /.config /.cache
    
    ENTRYPOINT ["yarn"]
    
    CMD ["build"]
    """
  end
end
