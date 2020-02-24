# -*- coding: utf-8 -*-
"""
Created on Sat Feb 15 22:32:18 2020

@author: Tanmay Thakur
"""

import pickle
import tensorflow as tf

from model import get_model
from tensorflow.keras.optimizers import Adam


X_train, y_train = pickle.load(open( "dict.pickle", "rb" ))

model = get_model(X_train)

model.compile(loss = 'mse', optimizer = Adam(lr = 1e-3))

cp_callbacks = tf.keras.callbacks.ModelCheckpoint(filepath = "recurrent_model_initial.h5", monitor = "val_loss", mode = 'min', save_best_only = True, verbose = 1)

model.fit(X_train, y_train, epochs = 20, batch_size = 16, validation_split = 0.25, callbacks = [cp_callbacks])
