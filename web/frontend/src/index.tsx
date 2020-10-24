import React from 'preact/compat';
import { h, render } from 'preact';
import 'preact/devtools';
import App from './App.js';
import './reset.css';
import './index.css';

const root: HTMLElement | null = document.getElementById('root');
if (root !== null) {
  render(<App />, root);
}
