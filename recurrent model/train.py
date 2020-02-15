# -*- coding: utf-8 -*-
"""
Created on Sat Feb 15 22:32:18 2020

@author: Tanmay Thakur
"""

import pickle

from model import get_model
from tensorflow.keras.optimizers import Adam


X_train, y_train = pickle.load(open( "dict.pickle", "rb" ))

model = get_model(X_train)

model.compile(loss = 'mse', optimizer = Adam(lr = 1e-3))

model.fit(X_train, y_train, epochs = 100, batch_size = 16, validation_split = 0.25)

model.save("recurrent_model_initial.h5")