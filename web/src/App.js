import React, { Component } from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route } from "react-router-dom";

import Navigation from './components/navigation';
import AlertList from './components/alertList';
import AddAlertForm from './components/addAlert'


class App extends Component {
  render() {
    return (
      <BrowserRouter>
      <div>
      <Navigation />
      <Route exact path="/" component={AlertList} />
      <Route path="/addalert" component={AddAlertForm} />
    </div>
    </BrowserRouter>
    );
  }
}

export default App;
