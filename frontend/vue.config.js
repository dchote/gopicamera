module.exports = {
  chainWebpack: config => {
    if (process.env.NODE_ENV === 'development') {
      config
        .output
        .filename('[name].[hash].js') 
        .end() 
    }  
  },
  runtimeCompiler: true,
  
  devServer: {
      host: '0.0.0.0',
      port: 8080,
      disableHostCheck: true,
      proxy: {
        '/config': {target: 'http://127.0.0.1:8000'},
        '/v1/*': {target: 'http://127.0.0.1:8000'},
        '/camera.mjpeg': {target: 'http://127.0.0.1:8000'},
      }
    }
}