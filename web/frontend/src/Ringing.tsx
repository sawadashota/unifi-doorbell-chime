import React, { useEffect, useState } from 'preact/compat';
// @ts-ignore
import { h, JSX } from 'preact';
import type { FunctionComponent } from 'preact';
import './Ringing.css';
import { Client } from './adapter/Client';

interface Prop {
  doorbell_id: string;
}

const Ringing: FunctionComponent<Prop> = ({ doorbell_id }) => {
  const [image, setImage] = useState<JSX.Element>(<p>loading snapshot...</p>);
  const [selected, select] = useState<string>('');
  const [templates, setTemplates] = useState<JSX.Element[]>([]);
  const noReactionText = 'No Reaction';

  const getTemplates = async (): Promise<void> => {
    const cl = await Client.configure();

    setImage(
      <img src={`${cl.apiEndpoint}/snapshot/${doorbell_id}`} alt="Snapshot" />,
    );

    const mt = await cl.messageTemplates();
    const els = mt.templates.map(
      (t: string, i: number): JSX.Element => {
        const onClick = () => {
          cl.setMessage(doorbell_id, t);
          select(t);
        };
        return (
          <button key={i} className="button blue block" onClick={onClick}>
            {t}
          </button>
        );
      },
    );
    setTemplates(els);
  };

  useEffect(() => {
    if (templates.length === 0) {
      getTemplates().catch(console.error);
    }
  }, []);

  const actions = () => {
    if (selected === '') {
      return (
        <div className="Ringing__buttons">
          {templates}
          <button
            className="button red block"
            onClick={() => select(noReactionText)}
          >
            {noReactionText}
          </button>
        </div>
      );
    }
    if (selected === noReactionText) {
      return (
        <div className="Ringing__buttons">
          <button className="button red block active">{selected}</button>
        </div>
      );
    }
    return (
      <div className="Ringing__buttons">
        <button className="button blue block active">{selected}</button>
      </div>
    );
  };
  return (
    <div className="Ringing">
      <h1 className="Ringing__title">Someone At The Door</h1>
      <div className="Ringing__snapshot">{image}</div>
      {actions()}
    </div>
  );
};

export { Ringing };
