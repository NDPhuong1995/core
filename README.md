# core
example beans:
  - registry bean: <br> &nbsp;&nbsp; beans.RegistryBean("defaultBean", defaultBean{})
  - get bean singleton: <br> &nbsp;&nbsp; bean := beans.GetBean("defaultBean", beans.ScopeSingleton)
  - get bean prototype: <br> &nbsp;&nbsp; bean := beans.GetBean("defaultBean", beans.Prototype)
  - invoke method bean: <br> &nbsp;&nbsp; beans.InvokeBeanMethod(bean, "method name", "args")
    
