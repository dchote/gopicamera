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
    }
}