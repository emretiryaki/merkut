import React, { Component } from 'react';
import Navigation from './components/navigation';
import Alert from './components/alertList';
import AddAlertForm from './components/addAlert'
class App extends Component {
  render() {
    return (
      <div>
      <Navigation />
      <Alert/>
    </div>
    );
  }
}

export default App;
