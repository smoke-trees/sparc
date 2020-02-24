# -*- coding: utf-8 -*-
"""
Created on Sat Feb 15 22:19:07 2020

@author: Tanmay Thakur
"""
import tensorflow as tf

from tensorflow.keras.models import Model
from tensorflow.keras.layers import Input, SpatialDropout1D, GRU, LSTM,Conv1D, concatenate, Dense
from tensorflow.keras.layers import GlobalAveragePooling1D, GlobalMaxPooling1D, Dropout


def get_model(X_train):
    inp = Input(shape=(X_train.shape[1],X_train.shape[2]))
    x = LSTM(512, activation = 'tanh', recurrent_activation = 'sigmoid', recurrent_dropout = 0, unroll = False, use_bias = True, return_sequences = True)(inp)
    x = LSTM(256, activation = 'tanh', recurrent_activation = 'sigmoid', recurrent_dropout = 0, unroll = False, use_bias = True, return_sequences = True)(inp)
    x = GlobalMaxPooling1D()(x)
    x = Dense(1024)(x)
    x = Dropout(0.25)(x)
    preds = Dense(X_train.shape[2])(x)
    
    model = Model(inp,preds)
    
    return model
