
var api = require ('../api/index')
var mocks = require('./mocks')

module.exports = (app) => {
    let obj = {};

    mocks.forEach (el => {
        if(el && el.key && api[el.key] && el.method){

            app[el.method](api[el.key]['url'] || '/issue/do',(req,res) => {
                
                res.json({
                    ...el.data
                })
            })
        }
    })

}

