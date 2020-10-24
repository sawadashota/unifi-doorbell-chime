import React from 'preact/compat';
import { h } from 'preact';
import Router from 'preact-router';
import AsyncRoute from 'preact-async-route';
import type { FunctionComponent } from 'preact';

const App: FunctionComponent = () => (
  <Router>
    <AsyncRoute
      path="/ringing/:doorbell_id"
      getComponent={() => import('./Ringing').then((module) => module.Ringing)}
    />
  </Router>
);

export default App;
